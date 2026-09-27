package column

import (
	"errors"
	"github.com/wyw14/solo-039/internal/model"
)

type Controller struct {
	State  *State
	Limits model.OperatingLimits
}

func NewController(state *State) *Controller {
	return &Controller{State: state, Limits: model.DefaultLimits()}
}

func (c *Controller) ApplySignals(campaign model.Campaign, pressure, level, temperature float64) error {
	if !model.PhaseAllowsHeat(campaign.Phase) {
		return errors.New("campaign does not accept heat signals")
	}
	c.State.Apply(campaign, pressure, level)
	c.State.mu.Lock()
	c.State.value.Temperature = temperature
	value := c.State.value
	c.State.mu.Unlock()
	return c.Limits.ValidateColumn(value)
}

func (c *Controller) Ready(campaign model.Campaign) bool {
	return PermitHeat(c.State, campaign) && c.State.Safe()
}

func (c *Controller) SnapshotValid() bool { return c.Limits.ValidateColumn(c.State.Snapshot()) == nil }
