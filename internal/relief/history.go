package relief

import (
	"github.com/wyw14/solo-039/internal/model"
	"sync"
	"time"
)

type Attempt struct {
	State model.Relief
	Error string
	At    time.Time
}

type History struct {
	mu       sync.Mutex
	attempts []Attempt
}

func NewHistory() *History { return &History{attempts: make([]Attempt, 0, 16)} }

func (h *History) Record(state model.Relief, err error) {
	value := Attempt{State: state, At: time.Now()}
	if err != nil {
		value.Error = err.Error()
	}
	h.mu.Lock()
	h.attempts = append(h.attempts, value)
	h.mu.Unlock()
}

func (h *History) Latest() (Attempt, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.attempts) == 0 {
		return Attempt{}, false
	}
	return h.attempts[len(h.attempts)-1], true
}

func (h *History) Failures() []Attempt {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Attempt, 0)
	for _, attempt := range h.attempts {
		if attempt.Error != "" || attempt.State.TimedOut {
			out = append(out, attempt)
		}
	}
	return out
}

func (h *History) Count() int { h.mu.Lock(); defer h.mu.Unlock(); return len(h.attempts) }
