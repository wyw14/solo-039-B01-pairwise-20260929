package utility

import (
	"github.com/wyw14/solo-039/internal/model"
	"sync"
)

type LeaseBook struct {
	mu     sync.Mutex
	leases map[string]model.Lease
}

func NewLeaseBook() *LeaseBook { return &LeaseBook{leases: make(map[string]model.Lease)} }

func (b *LeaseBook) Put(lease model.Lease) { b.mu.Lock(); b.leases[lease.Owner] = lease; b.mu.Unlock() }

func (b *LeaseBook) Get(owner string) (model.Lease, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	l, ok := b.leases[owner]
	return l, ok
}

func (b *LeaseBook) Fenced(owner string, fence uint64) bool {
	l, ok := b.Get(owner)
	return ok && l.Active && l.Fence == fence
}
