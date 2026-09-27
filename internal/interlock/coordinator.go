package interlock

import "github.com/wyw14/solo-039/internal/model"

type Coordinator struct{ Minimum int }

func NewCoordinator() *Coordinator { return &Coordinator{Minimum: 2} }

func (c *Coordinator) Evaluate(proofs []model.Proof, campaign model.Campaign) Permit {
	p := Permit{Evidence: append([]model.Proof(nil), proofs...)}
	if len(proofs) < c.Minimum {
		p.Reasons = append(p.Reasons, "insufficient evidence")
		return p
	}
	if !Independent(proofs) {
		p.Reasons = append(p.Reasons, "evidence is not independent")
		return p
	}
	if campaign.Phase != model.OnSpec {
		p.Reasons = append(p.Reasons, "campaign is not on-spec")
		return p
	}
	p.Allowed = true
	return p
}
