package service

import (
	"sync"
	"time"
)

type Metrics struct {
	mu         sync.Mutex
	operations map[string]int64
	failures   map[string]int64
	started    time.Time
	last       time.Time
}

func NewMetrics() *Metrics {
	return &Metrics{operations: make(map[string]int64), failures: make(map[string]int64), started: time.Now()}
}

func (m *Metrics) Record(name string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.operations[name]++
	if err != nil {
		m.failures[name]++
	}
	m.last = time.Now()
}

func (m *Metrics) Snapshot() map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()
	operations := make(map[string]int64, len(m.operations))
	failures := make(map[string]int64, len(m.failures))
	for key, value := range m.operations {
		operations[key] = value
	}
	for key, value := range m.failures {
		failures[key] = value
	}
	return map[string]any{"operations": operations, "failures": failures, "uptime_seconds": int64(time.Since(m.started).Seconds()), "last_operation_at": m.last.UnixNano()}
}

func (m *Metrics) FailureRate(name string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	total := m.operations[name]
	if total == 0 {
		return 0
	}
	return float64(m.failures[name]) / float64(total)
}
