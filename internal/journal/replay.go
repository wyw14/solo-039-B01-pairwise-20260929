package journal

import (
	"github.com/wyw14/solo-039/internal/model"
	"sort"
)

type Replay struct{ events []model.Event }

func NewReplay(events []model.Event) *Replay {
	copyEvents := append([]model.Event(nil), events...)
	sort.SliceStable(copyEvents, func(i, j int) bool { return copyEvents[i].At < copyEvents[j].At })
	return &Replay{events: copyEvents}
}

func (r *Replay) Campaign(id string, apply func(model.Event)) {
	for _, event := range r.events {
		if event.Campaign.String() == id {
			apply(event)
		}
	}
}

func (r *Replay) Events() []model.Event { return append([]model.Event(nil), r.events...) }
