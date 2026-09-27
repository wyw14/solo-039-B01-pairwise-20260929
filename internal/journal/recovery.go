package journal

import (
	"errors"
	"github.com/wyw14/solo-039/internal/model"
)

func ValidateRecovery(snapshot model.Snapshot, events []model.Event) error {
	if !model.ValidCampaign(snapshot.Campaign) {
		return errors.New("snapshot campaign is invalid")
	}
	if err := ValidateSequence(events); err != nil {
		return err
	}
	if latest, ok := Latest(events); ok && latest.Epoch > snapshot.Campaign.Epoch {
		return errors.New("snapshot trails durable events")
	}
	return nil
}

func RecoveryEvents(store *Store, campaign string) []model.Event {
	replay := NewReplay(store.Events())
	out := make([]model.Event, 0)
	replay.Campaign(campaign, func(event model.Event) { out = append(out, event) })
	return out
}
