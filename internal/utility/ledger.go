package utility

import (
	"github.com/wyw14/solo-039/internal/model"
	"sync"
)

type Ledger struct {
	mu       sync.Mutex
	leases   map[string]model.Lease
	sequence uint64
}

func NewLedger() *Ledger { return &Ledger{leases: make(map[string]model.Lease)} }

func (l *Ledger) Commit(lease model.Lease) uint64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sequence++
	lease.Fence = l.sequence
	lease.Active = true
	l.leases[lease.Owner] = lease
	return l.sequence
}

func (l *Ledger) Revoke(owner string, fence uint64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	lease, ok := l.leases[owner]
	if !ok || lease.Fence != fence {
		return false
	}
	lease.Active = false
	l.leases[owner] = lease
	return true
}

func (l *Ledger) Active() []model.Lease {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]model.Lease, 0, len(l.leases))
	for _, lease := range l.leases {
		if lease.Active {
			out = append(out, lease)
		}
	}
	return out
}

func (l *Ledger) Used() float64 {
	leases := l.Active()
	var total float64
	for _, lease := range leases {
		total += lease.Capacity
	}
	return total
}
