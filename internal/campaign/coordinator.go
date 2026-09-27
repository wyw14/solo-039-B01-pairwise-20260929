package campaign

import (
	"context"
	"errors"
	"github.com/wyw14/solo-039/internal/feed"
	"github.com/wyw14/solo-039/internal/model"
	"github.com/wyw14/solo-039/internal/reboiler"
	"github.com/wyw14/solo-039/internal/reflux"
)

type Coordinator struct {
	State    *State
	Feed     *feed.State
	Reflux   *reflux.State
	Reboiler *reboiler.State
}

func NewCoordinator(grade string) *Coordinator {
	return &Coordinator{State: New(grade), Feed: feed.New(), Reflux: reflux.New(), Reboiler: reboiler.New()}
}

func (c *Coordinator) Start(ctx context.Context) error {
	if err := c.Feed.Open(ctx, c.State.Snapshot()); err != nil {
		return err
	}
	c.State.Advance(model.Warming)
	if err := c.Reboiler.Ramp(ctx, c.State.Snapshot(), 0.45); err != nil {
		return err
	}
	c.State.Advance(model.Stabilizing)
	return nil
}

func (c *Coordinator) Cancel() { c.State.Stop(); c.Reflux.Stop(); c.Reboiler.Stop() }

func (c *Coordinator) EnsureRunning() error {
	if c.State.Snapshot().Phase == model.Stopped {
		return errors.New("campaign stopped")
	}
	return nil
}
