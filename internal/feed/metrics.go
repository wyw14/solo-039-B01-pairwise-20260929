package feed

import (
	"sync"
	"time"
)

type Metrics struct {
	mu            sync.Mutex
	switches      int
	activeSeconds int64
	lastSource    string
	lastAt        time.Time
}

func NewMetrics() *Metrics { return &Metrics{} }

func (m *Metrics) Record(source string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	if !m.lastAt.IsZero() {
		m.activeSeconds += int64(now.Sub(m.lastAt).Seconds())
	}
	m.lastSource = source
	m.lastAt = now
	m.switches++
}

func (m *Metrics) Snapshot() map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()
	return map[string]any{"switches": m.switches, "active_seconds": m.activeSeconds, "last_source": m.lastSource}
}

func (m *Metrics) Reset() {
	m.mu.Lock()
	m.switches = 0
	m.activeSeconds = 0
	m.lastSource = ""
	m.lastAt = time.Time{}
	m.mu.Unlock()
}
