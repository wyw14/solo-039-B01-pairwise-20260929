package journal

import "github.com/wyw14/solo-039/internal/model"

type RecoveredState struct {
	Campaign    model.Campaign
	LastEvent   model.Event
	Generations int
}

func RecoverState(events []model.Event, initial model.Campaign) RecoveredState {
	state := RecoveredState{Campaign: initial}
	generations := make(map[string]struct{})
	for _, event := range events {
		if event.Campaign != initial.ID {
			continue
		}
		generations[event.Revision.String()] = struct{}{}
		state.LastEvent = event
		state.Campaign.Revision = event.Revision
		state.Campaign.Epoch = event.Epoch
	}
	state.Generations = len(generations)
	return state
}

func (s RecoveredState) Ready() bool { return s.Campaign.ID != [16]byte{} && s.Campaign.Epoch > 0 }
