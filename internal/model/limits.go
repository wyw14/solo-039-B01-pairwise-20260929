package model

import (
	"errors"
	"math"
)

type OperatingLimits struct {
	MinPressure float64
	MaxPressure float64
	MinLevel    float64
	MaxLevel    float64
	MaxSteam    float64
	MaxReflux   float64
}

func DefaultLimits() OperatingLimits {
	return OperatingLimits{MinPressure: 0.2, MaxPressure: 2, MinLevel: 0.05, MaxLevel: 0.95, MaxSteam: 1, MaxReflux: 1}
}

func (l OperatingLimits) ValidateColumn(c Column) error {
	if math.IsNaN(c.Pressure) || math.IsNaN(c.Level) {
		return errors.New("column signal is not numeric")
	}
	if c.Pressure < l.MinPressure || c.Pressure > l.MaxPressure {
		return errors.New("column pressure outside operating limits")
	}
	if c.Level < l.MinLevel || c.Level > l.MaxLevel {
		return errors.New("column level outside operating limits")
	}
	return nil
}

func (l OperatingLimits) ValidateEnergy(reboiler Reboiler, reflux Reflux) error {
	if reboiler.Steam < 0 || reboiler.Steam > l.MaxSteam {
		return errors.New("steam outside operating limits")
	}
	if reflux.Ratio < 0 || reflux.Ratio > l.MaxReflux {
		return errors.New("reflux outside operating limits")
	}
	return nil
}

func (l OperatingLimits) Safe(snapshot Snapshot) bool {
	return l.ValidateColumn(snapshot.Column) == nil && l.ValidateEnergy(snapshot.Reboiler, snapshot.Reflux) == nil
}
