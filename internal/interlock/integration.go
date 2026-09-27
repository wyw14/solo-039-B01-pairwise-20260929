package interlock

import (
	"github.com/wyw14/solo-039/internal/column"
	"github.com/wyw14/solo-039/internal/model"
)

func HeatPermit(c *Coordinator, state *column.State, campaign model.Campaign) Permit {
	if !state.Safe() {
		return Permit{Reasons: []string{"column unsafe"}}
	}
	return c.Evaluate(nil, campaign)
}

func RoutePermit(c *Coordinator, proofs []model.Proof, campaign model.Campaign) bool {
	return c.Evaluate(proofs, campaign).Allowed
}
