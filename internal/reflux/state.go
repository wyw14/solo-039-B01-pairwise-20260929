package reflux

import (
	"context"
	"github.com/wyw14/solo-039/internal/model"
	"sync"
	"time"
)

type State struct {
	mu     sync.RWMutex
	value  model.Reflux
	cancel context.CancelFunc
}

func New() *State { return &State{value: model.Reflux{PumpCapacity: 1}} }

func (s *State) Snapshot() model.Reflux { s.mu.RLock(); defer s.mu.RUnlock(); return s.value }

func (s *State) Ramp(ctx context.Context, c model.Campaign, target float64, tick time.Duration) error {
	child, cancel := context.WithCancel(ctx)
	s.mu.Lock()
	s.cancel = cancel
	s.value.CommandEpoch = c.Epoch
	s.mu.Unlock()
	defer cancel()
	for {
		s.mu.Lock()
		if s.value.Ratio >= target {
			s.mu.Unlock()
			return nil
		}
		s.value.Ratio += 0.1
		s.mu.Unlock()
		select {
		case <-child.Done():
			return child.Err()
		case <-time.After(tick):
		}
	}
}

func (s *State) Stop() {
	s.mu.Lock()
	cancel := s.cancel
	s.cancel = nil
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (s *State) SetLease(lease model.Lease) { s.mu.Lock(); s.value.Lease = lease; s.mu.Unlock() }
