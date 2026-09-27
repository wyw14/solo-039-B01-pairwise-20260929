package reflux

import (
	"context"
	"github.com/wyw14/solo-039/internal/model"
	"time"
)

type Controller struct {
	State *State
	Tick  time.Duration
}

func NewController() *Controller { return &Controller{State: New(), Tick: time.Millisecond} }

func (c *Controller) Ramp(ctx context.Context, campaign model.Campaign, target float64) error {
	return c.State.Ramp(ctx, campaign, target, c.Tick)
}

func (c *Controller) Stop() { c.State.Stop() }
