package relief

import (
	"errors"
	"github.com/wyw14/solo-039/internal/model"
)

type Evaluation struct {
	PressureSafe  bool
	ValveReleased bool
	TimedOut      bool
	Message       string
}

func Evaluate(state model.Relief, safePressure float64) Evaluation {
	value := Evaluation{PressureSafe: state.Pressure <= safePressure, ValveReleased: state.Released, TimedOut: state.TimedOut}
	switch {
	case value.TimedOut:
		value.Message = "relief valve timed out"
	case !value.PressureSafe:
		value.Message = "pressure remains above safe threshold"
	case !value.ValveReleased:
		value.Message = "relief path not confirmed"
	default:
		value.Message = "depressurization complete"
	}
	return value
}

func (e Evaluation) Error() error {
	if e.PressureSafe && e.ValveReleased && !e.TimedOut {
		return nil
	}
	return errors.New(e.Message)
}

func (e Evaluation) Safe() bool { return e.Error() == nil }
