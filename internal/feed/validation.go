package feed

import (
	"errors"
	"github.com/google/uuid"
)

type SwitchProof struct {
	OldRevision uuid.UUID
	NewRevision uuid.UUID
	OldClosed   bool
	OldFlow     float64
	NewOpened   bool
}

func BuildSwitchProof(old, next *Route) SwitchProof {
	proof := SwitchProof{}
	if old != nil {
		proof.OldRevision = old.Revision
		proof.OldClosed = !old.Open
		proof.OldFlow = old.Flow
	}
	if next != nil {
		proof.NewRevision = next.Revision
		proof.NewOpened = next.Open
	}
	return proof
}

func (p SwitchProof) Validate() error {
	if p.OldRevision == uuid.Nil || p.NewRevision == uuid.Nil {
		return errors.New("feed revisions are incomplete")
	}
	if p.OldRevision == p.NewRevision {
		return errors.New("feed revision did not advance")
	}
	if !p.OldClosed || p.OldFlow != 0 {
		return errors.New("old feed source remains active")
	}
	if !p.NewOpened {
		return errors.New("new feed source is not active")
	}
	return nil
}

func (p SwitchProof) Safe() bool { return p.Validate() == nil }

func (c *Coordinator) Proof() SwitchProof { return BuildSwitchProof(c.Previous, c.Current) }
