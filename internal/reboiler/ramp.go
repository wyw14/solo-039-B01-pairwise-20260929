package reboiler

import (
	"context"
	"github.com/wyw14/solo-039/internal/model"
)

type Ramp struct{ State *State }

func NewRamp() *Ramp { return &Ramp{State: New()} }

func (r *Ramp) Apply(ctx context.Context, c model.Campaign, target float64) error {
	return r.State.Ramp(ctx, c, target)
}
