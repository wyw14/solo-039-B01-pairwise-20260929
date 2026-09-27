package service

import (
	"context"
	"github.com/wyw14/solo-039/internal/campaign"
	"github.com/wyw14/solo-039/internal/model"
)

type Coordinator struct{ Campaign *campaign.Coordinator }

func NewCoordinator(grade string) *Coordinator {
	return &Coordinator{Campaign: campaign.NewCoordinator(grade)}
}

func (c *Coordinator) Run(ctx context.Context) error { return c.Campaign.Start(ctx) }

func (c *Coordinator) Cancel() { c.Campaign.Cancel() }

func (c *Coordinator) Phase() model.Phase { return c.Campaign.State.Snapshot().Phase }
