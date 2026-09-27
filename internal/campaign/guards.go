package campaign

import (
	"errors"
	"github.com/wyw14/solo-039/internal/model"
)

func AdvanceChecked(state *State, next model.Phase) (model.Campaign, error) {
	current := state.Snapshot()
	if err := model.TransitionError(current.Phase, next); err != nil {
		return current, err
	}
	return state.Advance(next), nil
}

func RequirePhase(state *State, expected ...model.Phase) error {
	current := state.Snapshot()
	for _, phase := range expected {
		if current.Phase == phase {
			return nil
		}
	}
	return errors.New("campaign phase is not permitted")
}

func RevisionMatches(state *State, revision string) bool {
	return state.Snapshot().Revision.String() == revision
}

func Restart(state *State) model.Campaign {
	current := state.Snapshot()
	current.Phase = model.Preparing
	current.Epoch++
	current.Fencing++
	Restore(state, current)
	return current
}
