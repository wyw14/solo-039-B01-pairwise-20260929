package column

import "github.com/wyw14/solo-039/internal/model"

func Stabilize(s *State, c model.Campaign, pressure, level float64) bool {
	if !model.PhaseAllowsHeat(c.Phase) {
		return false
	}
	s.Apply(c, pressure, level)
	return s.Safe()
}
