package reboiler

import (
	"context"
	"github.com/wyw14/solo-039/internal/model"
	"sync"
	"time"
)

type State struct {
	mu       sync.RWMutex
	value    model.Reboiler
	cancel   context.CancelFunc
	ackEpoch uint64
}

func New() *State { return &State{value: model.Reboiler{Closed: true}} }

func (s *State) Snapshot() model.Reboiler { s.mu.RLock(); defer s.mu.RUnlock(); return s.value }

func (s *State) Ramp(ctx context.Context, c model.Campaign, target float64) error {
	child, cancel := context.WithCancel(ctx)
	s.mu.Lock()
	s.cancel = cancel
	s.value.CommandEpoch = c.Epoch
	s.value.Closed = false
	s.mu.Unlock()
	defer cancel()
	for {
		s.mu.Lock()
		if s.value.Steam >= target {
			s.mu.Unlock()
			return nil
		}
		s.value.Steam += 0.1
		s.mu.Unlock()
		select {
		case <-child.Done():
			return child.Err()
		case <-time.After(time.Millisecond):
		}
	}
}

func (s *State) Close(ctx context.Context, epoch uint64) error {
	s.mu.Lock()
	s.value.CommandEpoch = epoch
	s.value.Steam = 0
	s.value.Closed = false
	s.value.Trip++
	s.mu.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s *State) AcknowledgeClosed(epoch uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if epoch != s.value.CommandEpoch {
		return false
	}
	s.ackEpoch = epoch
	s.value.Closed = true
	return true
}

func (s *State) ClosureConfirmed(epoch uint64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value.Closed && s.ackEpoch == epoch && s.value.CommandEpoch == epoch
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
