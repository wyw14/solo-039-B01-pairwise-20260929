package journal

import (
	"errors"
	"github.com/wyw14/solo-039/internal/model"
	"sync"
)

type Store struct {
	mu       sync.Mutex
	events   []model.Event
	failNext bool
}

func New() *Store { return &Store{events: make([]model.Event, 0, 64)} }

func (s *Store) FailNext() { s.mu.Lock(); s.failNext = true; s.mu.Unlock() }

func (s *Store) Append(event model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.failNext {
		s.failNext = false
		return errors.New("journal short write")
	}
	s.events = append(s.events, event)
	return nil
}

func (s *Store) Events() []model.Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	copyEvents := make([]model.Event, len(s.events))
	copy(copyEvents, s.events)
	return copyEvents
}

func (s *Store) Last(kind string) (model.Event, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := len(s.events) - 1; i >= 0; i-- {
		if s.events[i].Kind == kind {
			return s.events[i], true
		}
	}
	return model.Event{}, false
}
