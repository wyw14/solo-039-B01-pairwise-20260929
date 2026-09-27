package service

import (
	"github.com/google/uuid"
	"sort"
	"sync"
	"time"
)

type Incident struct {
	ID       uuid.UUID
	Kind     string
	Message  string
	Revision string
	Open     bool
	At       time.Time
}

type Incidents struct {
	mu    sync.Mutex
	items map[uuid.UUID]Incident
}

func NewIncidents() *Incidents { return &Incidents{items: make(map[uuid.UUID]Incident)} }

func (s *Incidents) Open(kind, message, revision string) Incident {
	item := Incident{ID: uuid.New(), Kind: kind, Message: message, Revision: revision, Open: true, At: time.Now()}
	s.mu.Lock()
	s.items[item.ID] = item
	s.mu.Unlock()
	return item
}

func (s *Incidents) Resolve(id uuid.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok {
		return false
	}
	item.Open = false
	s.items[id] = item
	return true
}

func (s *Incidents) List(openOnly bool) []Incident {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Incident, 0, len(s.items))
	for _, item := range s.items {
		if openOnly && !item.Open {
			continue
		}
		out = append(out, item)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}

func (s *Incidents) CountOpen() int { return len(s.List(true)) }
