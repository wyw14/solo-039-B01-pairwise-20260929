package reboiler

import (
	"context"
	"github.com/wyw14/solo-039/internal/model"
)

func CloseForTrip(ctx context.Context, s *State, c model.Campaign) error {
	return s.Close(ctx, c.Epoch)
}

func Permit(s *State, c model.Campaign) bool {
	return s.ClosureConfirmed(c.Epoch)
}
