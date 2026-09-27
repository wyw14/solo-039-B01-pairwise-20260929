package journal

import (
	"errors"
	"github.com/wyw14/solo-039/internal/model"
	"sort"
)

type RecoveryStep struct {
	Name     string
	Requires []string
	Ready    bool
	Applied  bool
}

type RecoveryPlan struct {
	steps map[string]*RecoveryStep
	order []string
}

func NewRecoveryPlan() *RecoveryPlan { return &RecoveryPlan{steps: make(map[string]*RecoveryStep)} }

func (p *RecoveryPlan) Add(name string, requires ...string) error {
	if name == "" {
		return errors.New("recovery step name required")
	}
	if _, exists := p.steps[name]; exists {
		return errors.New("duplicate recovery step")
	}
	step := &RecoveryStep{Name: name, Requires: append([]string(nil), requires...)}
	p.steps[name] = step
	p.order = append(p.order, name)
	return nil
}

func (p *RecoveryPlan) MarkReady(name string) bool {
	step := p.steps[name]
	if step == nil {
		return false
	}
	for _, dependency := range step.Requires {
		required := p.steps[dependency]
		if required == nil || !required.Applied {
			return false
		}
	}
	step.Ready = true
	return true
}

func (p *RecoveryPlan) Apply(name string) bool {
	step := p.steps[name]
	if step == nil || !step.Ready {
		return false
	}
	step.Applied = true
	return true
}

func (p *RecoveryPlan) Pending() []string {
	pending := make([]string, 0)
	for _, name := range p.order {
		if !p.steps[name].Applied {
			pending = append(pending, name)
		}
	}
	return pending
}

func (p *RecoveryPlan) Complete() bool { return len(p.Pending()) == 0 }

func DefaultRecoveryPlan() *RecoveryPlan {
	plan := NewRecoveryPlan()
	_ = plan.Add("campaign")
	_ = plan.Add("column", "campaign")
	_ = plan.Add("fencing", "campaign", "column")
	_ = plan.Add("interlocks", "fencing")
	_ = plan.Add("controllers", "interlocks")
	return plan
}

func ReplayByGeneration(events []model.Event) map[string][]model.Event {
	groups := make(map[string][]model.Event)
	for _, event := range events {
		key := event.Revision.String()
		groups[key] = append(groups[key], event)
	}
	for key := range groups {
		sort.SliceStable(groups[key], func(i, j int) bool { return groups[key][i].At < groups[key][j].At })
	}
	return groups
}
