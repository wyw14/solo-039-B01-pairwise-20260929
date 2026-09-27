package service

import (
	"github.com/wyw14/solo-039/internal/campaign"
	"github.com/wyw14/solo-039/internal/column"
	"github.com/wyw14/solo-039/internal/journal"
	"github.com/wyw14/solo-039/internal/model"
	"github.com/wyw14/solo-039/internal/reboiler"
)

func Recover(store *journal.Store, runtime *Runtime, boiler *reboiler.State) error {
	events := journal.NewReplay(store.Events())
	var snapshot model.Campaign
	events.Campaign(runtime.Campaign.Snapshot().ID.String(), func(event model.Event) { snapshot = runtime.Campaign.Snapshot(); snapshot.Epoch = event.Epoch })
	if snapshot.ID != [16]byte{} {
		column.Restore(runtime.Column, model.Column{ID: "T-101", Phase: snapshot.Phase, Fencing: snapshot.Fencing, Pressure: 0.8, Level: 0.4})
		campaign.Restore(runtime.Campaign, snapshot)
	}
	if !model.PhaseAllowsHeat(runtime.Campaign.Snapshot().Phase) {
		boiler.Stop()
	}
	return nil
}
