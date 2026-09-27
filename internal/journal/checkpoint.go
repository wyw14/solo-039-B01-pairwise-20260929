package journal

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/wyw14/solo-039/internal/model"
)

type Checkpoint struct {
	EventID  string
	Revision string
	Epoch    uint64
	Digest   string
}

func NewCheckpoint(event model.Event) Checkpoint {
	raw := event.Kind + "|" + event.Payload + "|" + event.Revision.String()
	hash := sha256.Sum256([]byte(raw))
	return Checkpoint{EventID: event.ID.String(), Revision: event.Revision.String(), Epoch: event.Epoch, Digest: hex.EncodeToString(hash[:])}
}

func (c Checkpoint) Valid(event model.Event) bool {
	return c.EventID == event.ID.String() && c.Revision == event.Revision.String() && c.Epoch == event.Epoch && c.Digest != ""
}

func ValidateSequence(events []model.Event) error {
	var epoch uint64
	for _, event := range events {
		if event.Epoch < epoch {
			return errors.New("event epoch regressed")
		}
		epoch = event.Epoch
	}
	return nil
}

func Latest(events []model.Event) (model.Event, bool) {
	if len(events) == 0 {
		return model.Event{}, false
	}
	return events[len(events)-1], true
}
