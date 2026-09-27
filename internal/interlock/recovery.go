package interlock

import (
	"github.com/wyw14/solo-039/internal/model"
	"github.com/wyw14/solo-039/internal/utility"
)

type RecoveryState struct {
	Campaign model.Campaign
	Lease    model.Lease
	Proofs   []model.Proof
}

func RestorePermit(coordinator *Coordinator, state RecoveryState) Permit {
	if !utility.ValidLease(state.Lease, state.Campaign.Fencing) {
		return Permit{Reasons: []string{"restored lease is fenced"}, Evidence: state.Proofs}
	}
	return coordinator.Evaluate(state.Proofs, state.Campaign)
}

func RecoverySafe(coordinator *Coordinator, state RecoveryState) bool {
	return RestorePermit(coordinator, state).Allowed
}
