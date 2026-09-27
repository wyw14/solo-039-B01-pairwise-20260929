package column

import "github.com/wyw14/solo-039/internal/model"

func Restore(s *State, snapshot model.Column) { s.mu.Lock(); s.value = snapshot; s.mu.Unlock() }

func PermitHeat(s *State, campaign model.Campaign) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return model.PhaseAllowsHeat(campaign.Phase) && model.PhaseAllowsHeat(s.value.Phase) && campaign.Fencing >= s.value.Fencing
}

func ColumnReady(c model.Column) bool { return model.PhaseAllowsHeat(c.Phase) && c.Pressure < 2 }
