package model

import (
	"math"
	"strings"
	"testing"
)

func validEnergy() (Reboiler, Reflux) {
	return Reboiler{Steam: 0.5}, Reflux{Ratio: 0.5}
}

func TestValidateEnergyRejectsNonFinite(t *testing.T) {
	limits := DefaultLimits()
	badNumbers := map[string]float64{
		"nan":      math.NaN(),
		"+inf":     math.Inf(1),
		"-inf":     math.Inf(-1),
	}

	for name, bad := range badNumbers {
		reboiler, reflux := validEnergy()
		reboiler.Steam = bad
		err := limits.ValidateEnergy(reboiler, reflux)
		if err == nil || !strings.Contains(err.Error(), "steam") {
			t.Errorf("steam %s: expected steam signal error, got %v", name, err)
		}
		if limits.Safe(Snapshot{Column: validColumn(), Reboiler: reboiler, Reflux: reflux}) {
			t.Errorf("steam %s: snapshot must not be safe", name)
		}

		reboiler, reflux = validEnergy()
		reflux.Ratio = bad
		err = limits.ValidateEnergy(reboiler, reflux)
		if err == nil || !strings.Contains(err.Error(), "reflux") {
			t.Errorf("reflux %s: expected reflux signal error, got %v", name, err)
		}
		if limits.Safe(Snapshot{Column: validColumn(), Reboiler: reboiler, Reflux: reflux}) {
			t.Errorf("reflux %s: snapshot must not be safe", name)
		}
	}
}

func TestValidateEnergyBoundaries(t *testing.T) {
	limits := DefaultLimits()
	for _, steam := range []float64{0, limits.MaxSteam} {
		for _, ratio := range []float64{0, limits.MaxReflux} {
			reboiler, reflux := Reboiler{Steam: steam}, Reflux{Ratio: ratio}
			if err := limits.ValidateEnergy(reboiler, reflux); err != nil {
				t.Errorf("steam=%v ratio=%v: expected safe, got %v", steam, ratio, err)
			}
		}
	}

	reboiler, reflux := validEnergy()
	reboiler.Steam = limits.MaxSteam + 0.01
	if err := limits.ValidateEnergy(reboiler, reflux); err == nil {
		t.Error("steam above max must be rejected")
	}

	reboiler, reflux = validEnergy()
	reflux.Ratio = -0.01
	if err := limits.ValidateEnergy(reboiler, reflux); err == nil {
		t.Error("negative reflux must be rejected")
	}
}

func validColumn() Column {
	return Column{
		Pressure: (DefaultLimits().MinPressure + DefaultLimits().MaxPressure) / 2,
		Level:    (DefaultLimits().MinLevel + DefaultLimits().MaxLevel) / 2,
	}
}
