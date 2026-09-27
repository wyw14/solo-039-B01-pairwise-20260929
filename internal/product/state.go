package product

import (
	"github.com/wyw14/solo-039/internal/model"
	"sync"
)

type State struct {
	mu     sync.RWMutex
	routes map[string]Route
}

func NewState() *State { return &State{routes: make(map[string]Route)} }

func (s *State) Put(plan model.ProductPlan, route Route) {
	s.mu.Lock()
	s.routes[plan.Target] = route
	s.mu.Unlock()
}

func (s *State) Get(target string) (Route, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.routes[target]
	return r, ok
}

func (s *State) Targets() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.routes))
	for target := range s.routes {
		out = append(out, target)
	}
	return out
}
