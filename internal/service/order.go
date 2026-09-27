package service

import (
	"errors"
	"github.com/wyw14/solo-039/internal/journal"
)

type RecoveryHooks struct {
	LoadCampaign     func() error
	LoadColumn       func() error
	LoadFencing      func() error
	LoadInterlocks   func() error
	StartControllers func() error
}

func RecoverOrdered(hooks RecoveryHooks) ([]string, error) {
	plan := journal.DefaultRecoveryPlan()
	steps := []struct {
		name string
		run  func() error
	}{
		{name: "campaign", run: hooks.LoadCampaign},
		{name: "column", run: hooks.LoadColumn},
		{name: "fencing", run: hooks.LoadFencing},
		{name: "interlocks", run: hooks.LoadInterlocks},
		{name: "controllers", run: hooks.StartControllers},
	}
	applied := make([]string, 0, len(steps))
	for _, step := range steps {
		if step.run == nil {
			return applied, errors.New("recovery hook is missing")
		}
		if !plan.MarkReady(step.name) {
			return applied, errors.New("recovery dependency is not ready")
		}
		if err := step.run(); err != nil {
			return applied, err
		}
		if !plan.Apply(step.name) {
			return applied, errors.New("recovery step could not be applied")
		}
		applied = append(applied, step.name)
	}
	if !plan.Complete() {
		return applied, errors.New("recovery did not complete")
	}
	return applied, nil
}
