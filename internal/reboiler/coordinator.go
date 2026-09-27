package reboiler

import (
	"context"
	"github.com/wyw14/solo-039/internal/model"
)

type Coordinator struct{ State *State }

func NewCoordinator() *Coordinator { return &Coordinator{State: New()} }

func (c *Coordinator) Ramp(ctx context.Context, campaign model.Campaign, target float64) error {
	return c.State.Ramp(ctx, campaign, target)
}

func (c *Coordinator) Stop() { c.State.Stop() }
