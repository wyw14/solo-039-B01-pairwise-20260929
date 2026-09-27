package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
)

func NewCampaign(grade string) Campaign {
	return Campaign{ID: uuid.New(), Grade: grade, Revision: uuid.New(), Phase: Preparing, Epoch: 1, Fencing: 1}
}

func NewEvent(kind string, c Campaign, payload string) Event {
	return Event{ID: uuid.New(), Kind: kind, Campaign: c.ID, Revision: c.Revision, Epoch: c.Epoch, Payload: payload}
}

func ProofDigest(p Proof) string {
	value := fmt.Sprintf("%s|%s|%s|%d", p.ID, p.Revision, p.Source, p.At)
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func ValidCampaign(c Campaign) bool {
	return c.ID != uuid.Nil && c.Revision != uuid.Nil && c.Epoch > 0 && c.Fencing > 0
}
