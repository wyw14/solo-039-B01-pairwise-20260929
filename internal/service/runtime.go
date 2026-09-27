package service

import (
	"context"
	"github.com/wyw14/solo-039/internal/campaign"
	"github.com/wyw14/solo-039/internal/column"
	"github.com/wyw14/solo-039/internal/model"
)

type Runtime struct {
	Campaign *campaign.State
	Column   *column.State
	running  bool
}

func NewRuntime(grade string) *Runtime {
	return &Runtime{Campaign: campaign.New(grade), Column: column.New("T-101")}
}

func (r *Runtime) Start(ctx context.Context) error {
	r.Campaign.Advance(model.Warming)
	r.Column.Apply(r.Campaign.Snapshot(), 1.2, 0.5)
	r.running = true
	return nil
}

func (r *Runtime) Running() bool { return r.running }

func (r *Runtime) Stop() { r.Campaign.Stop(); r.running = false }
