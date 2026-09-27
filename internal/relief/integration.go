package relief

import (
	"context"
	"github.com/wyw14/solo-039/internal/model"
	"github.com/wyw14/solo-039/internal/reboiler"
)

func Isolate(ctx context.Context, state *State, boiler *reboiler.State, c model.Campaign) error {
	if err := boiler.Close(ctx, c.Epoch); err != nil {
		return err
	}
	if state.Snapshot().TimedOut {
		return context.DeadlineExceeded
	}
	state.Release()
	return nil
}
