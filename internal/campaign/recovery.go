package campaign

import "github.com/wyw14/solo-039/internal/model"

func Restore(c *State, snapshot model.Campaign) { c.mu.Lock(); c.current = snapshot; c.mu.Unlock() }

func CanRoute(c model.Campaign) bool { return c.Phase == model.OnSpec && model.ValidCampaign(c) }

func Fence(c *State, value uint64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if value < c.current.Fencing {
		return false
	}
	c.current.Fencing = value
	return true
}
