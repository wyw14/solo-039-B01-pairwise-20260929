package reflux

import (
	"github.com/wyw14/solo-039/internal/model"
	"sync"
	"time"
)

type Telemetry struct {
	mu      sync.Mutex
	samples []float64
	updated time.Time
}

func NewTelemetry() *Telemetry { return &Telemetry{samples: make([]float64, 0, 32)} }

func (t *Telemetry) Record(state model.Reflux) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.samples = append(t.samples, state.Ratio)
	if len(t.samples) > 32 {
		t.samples = t.samples[len(t.samples)-32:]
	}
	t.updated = time.Now()
}

func (t *Telemetry) Mean() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.samples) == 0 {
		return 0
	}
	var sum float64
	for _, value := range t.samples {
		sum += value
	}
	return sum / float64(len(t.samples))
}

func (t *Telemetry) Stable(tolerance float64) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.samples) < 2 {
		return false
	}
	last := t.samples[len(t.samples)-1]
	for _, value := range t.samples[:len(t.samples)-1] {
		if value-last > tolerance || last-value > tolerance {
			return false
		}
	}
	return true
}

func (t *Telemetry) Snapshot() map[string]any {
	return map[string]any{"mean": t.Mean(), "stable": t.Stable(0.05), "updated": t.updated.UnixNano()}
}
