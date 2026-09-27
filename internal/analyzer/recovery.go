package analyzer

import (
	"github.com/google/uuid"
	"github.com/wyw14/solo-039/internal/model"
)

type RecoveryState struct {
	Calibration  uuid.UUID
	FeedRevision uuid.UUID
	Samples      []model.Sample
}

func Capture(state *State) RecoveryState {
	calibration, feedRevision := state.window.Identity()
	return RecoveryState{Calibration: calibration, FeedRevision: feedRevision, Samples: state.window.Samples()}
}

func Restore(state *State, recovery RecoveryState) {
	state.mu.Lock()
	state.calibration = recovery.Calibration
	state.window.Reset(recovery.Calibration, recovery.FeedRevision)
	state.mu.Unlock()
	for _, sample := range recovery.Samples {
		state.Add(sample)
	}
}

func Compatible(recovery RecoveryState, calibration, feedRevision uuid.UUID) bool {
	return recovery.Calibration == calibration && recovery.FeedRevision == feedRevision
}
