package utility

import (
	"errors"
	"github.com/wyw14/solo-039/internal/model"
	"sync"
)

type Capacity struct {
	mu    sync.Mutex
	total float64
	used  float64
}

func NewCapacity(total float64) *Capacity { return &Capacity{total: total} }

func (c *Capacity) Reserve(owner string, amount float64, fence uint64) (model.Lease, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if amount <= 0 || c.used+amount > c.total {
		return model.Lease{}, errors.New("capacity unavailable")
	}
	c.used += amount
	return model.Lease{Owner: owner, Capacity: amount, Fence: fence, Active: true}, nil
}

func (c *Capacity) Release(lease model.Lease) {
	c.mu.Lock()
	if lease.Active {
		c.used -= lease.Capacity
	}
	c.mu.Unlock()
}

func (c *Capacity) Used() float64 { c.mu.Lock(); defer c.mu.Unlock(); return c.used }
