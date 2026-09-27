package reboiler

import (
	"github.com/wyw14/solo-039/internal/model"
	"sync"
	"time"
)

type Telemetry struct {
	mu      sync.Mutex
	values  []float64
	trips   int
	updated time.Time
}

func NewTelemetry() *Telemetry { return &Telemetry{values: make([]float64, 0, 32)} }

func (t *Telemetry) Record(state model.Reboiler) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.values = append(t.values, state.Steam)
	if state.Trip > uint64(t.trips) {
		t.trips = int(state.Trip)
	}
	if len(t.values) > 32 {
		t.values = t.values[len(t.values)-32:]
	}
	t.updated = time.Now()
}

func (t *Telemetry) Last() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.values) == 0 {
		return 0
	}
	return t.values[len(t.values)-1]
}

func (t *Telemetry) Trips() int { t.mu.Lock(); defer t.mu.Unlock(); return t.trips }

func (t *Telemetry) Snapshot() map[string]any {
	return map[string]any{"steam": t.Last(), "trips": t.Trips(), "updated": t.updated.UnixNano()}
}
