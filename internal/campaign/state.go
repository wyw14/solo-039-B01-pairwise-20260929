package campaign

import (
	"context"
	"github.com/wyw14/solo-039/internal/model"
	"sync"
)

type State struct {
	mu      sync.RWMutex
	current model.Campaign
	cancel  context.CancelFunc
}

func New(grade string) *State { return &State{current: model.NewCampaign(grade)} }

func (s *State) Snapshot() model.Campaign { s.mu.RLock(); defer s.mu.RUnlock(); return s.current }

func (s *State) Advance(phase model.Phase) model.Campaign {
	s.mu.Lock()
	s.current.Phase = phase
	s.current.Epoch++
	s.mu.Unlock()
	return s.Snapshot()
}

func (s *State) Stop() {
	s.mu.Lock()
	s.current.Phase = model.Stopped
	s.current.Epoch++
	cancel := s.cancel
	s.cancel = nil
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (s *State) Bind(ctx context.Context) context.Context {
	child, cancel := context.WithCancel(ctx)
	s.mu.Lock()
	s.cancel = cancel
	s.mu.Unlock()
	return child
}
