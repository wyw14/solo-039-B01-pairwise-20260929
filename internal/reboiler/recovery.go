package reboiler

import (
	"errors"
	"github.com/wyw14/solo-039/internal/model"
)

func Restore(state *State, snapshot model.Reboiler, campaign model.Campaign) error {
	if snapshot.CommandEpoch > campaign.Epoch {
		return errors.New("reboiler snapshot is ahead of campaign")
	}
	state.mu.Lock()
	state.value = snapshot
	if snapshot.Closed {
		state.ackEpoch = snapshot.CommandEpoch
	}
	state.mu.Unlock()
	return nil
}

func Recoverable(snapshot model.Reboiler, campaign model.Campaign) bool {
	return snapshot.CommandEpoch <= campaign.Epoch && (model.PhaseAllowsHeat(campaign.Phase) || snapshot.Closed)
}
