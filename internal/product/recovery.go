package product

import (
	"errors"
	"github.com/wyw14/solo-039/internal/journal"
	"github.com/wyw14/solo-039/internal/model"
)

type Recovery struct {
	Events    []model.Event
	Decisions []Decision
}

func BuildRecovery(events []model.Event, audit *Audit) Recovery {
	decisions := make([]Decision, 0)
	if latest, err := audit.Latest(); err == nil {
		decisions = append(decisions, latest)
	}
	return Recovery{Events: append([]model.Event(nil), events...), Decisions: decisions}
}

func (r Recovery) Durable(revision string) bool {
	for _, event := range r.Events {
		if event.Kind == "proof.durable" && event.Revision.String() == revision {
			return true
		}
	}
	return false
}

func (r Recovery) Routed(revision string) bool {
	for _, decision := range r.Decisions {
		if decision.Opened && decision.Plan.Revision.String() == revision {
			return true
		}
	}
	return false
}

func (r Recovery) Validate(revision string) error {
	if r.Routed(revision) && !r.Durable(revision) {
		return errors.New("routed product has no durable proof")
	}
	return nil
}

func RestoreRoute(coordinator *Coordinator, plan model.ProductPlan, store *journal.Store, recovery Recovery) error {
	if err := recovery.Validate(plan.Revision.String()); err != nil {
		return err
	}
	if !recovery.Durable(plan.Revision.String()) {
		return errors.New("product proof unavailable during recovery")
	}
	return coordinator.Recover(plan)
}
