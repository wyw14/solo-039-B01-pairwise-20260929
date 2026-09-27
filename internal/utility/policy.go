package utility

import "github.com/wyw14/solo-039/internal/model"

func ValidLease(lease model.Lease, fence uint64) bool {
	return lease.Active && lease.Capacity > 0 && lease.Fence == fence
}

func CanOpen(lease model.Lease, owner string) bool { return lease.Active && lease.Owner == owner }
