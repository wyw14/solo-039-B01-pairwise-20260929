package campaign

import (
	"github.com/wyw14/solo-039/internal/model"
	"github.com/wyw14/solo-039/internal/utility"
)

func ApplyLease(state *State, lease model.Lease) bool {
	if !utility.ValidLease(lease, state.Snapshot().Fencing) {
		return false
	}
	return Fence(state, lease.Fence)
}

func TransitionRevision(state *State, grade string) model.Campaign {
	current := state.Advance(model.Transitioning)
	current.Grade = grade
	current.Revision = model.NewCampaign(grade).Revision
	Restore(state, current)
	return current
}

func CompleteTransition(state *State) model.Campaign { return state.Advance(model.Stabilizing) }
