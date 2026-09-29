package model

import (
	"math"
	"strings"
	"testing"
)

func TestValidateEnergyRejectsNonFinite(t *testing.T) {
	limits := DefaultLimits()

	cases := []struct {
		name      string
		steam     float64
		ratio     float64
		wantInErr string
	}{
		{"steam NaN", math.NaN(), 0.5, "steam"},
		{"steam +Inf", math.Inf(1), 0.5, "steam"},
		{"steam -Inf", math.Inf(-1), 0.5, "steam"},
		{"reflux NaN", 0.5, math.NaN(), "reflux"},
		{"reflux +Inf", 0.5, math.Inf(1), "reflux"},
		{"reflux -Inf", 0.5, math.Inf(-1), "reflux"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := limits.ValidateEnergy(Reboiler{Steam: tc.steam}, Reflux{Ratio: tc.ratio})
			if err == nil {
				t.Fatalf("expected error for %s, got nil", tc.name)
			}
			if !strings.Contains(err.Error(), tc.wantInErr) {
				t.Fatalf("error %q does not identify the %q signal", err.Error(), tc.wantInErr)
			}
			snapshot := Snapshot{
				Column:   Column{Pressure: 1, Level: 0.5},
				Reboiler: Reboiler{Steam: tc.steam},
				Reflux:   Reflux{Ratio: tc.ratio},
			}
			if limits.Safe(snapshot) {
				t.Fatalf("snapshot with %s was judged safe", tc.name)
			}
		})
	}
}

func TestValidateEnergyAcceptsBoundaries(t *testing.T) {
	limits := DefaultLimits()

	cases := []struct {
		name  string
		steam float64
		ratio float64
	}{
		{"lower boundaries", 0, 0},
		{"upper boundaries", limits.MaxSteam, limits.MaxReflux},
		{"mid range", limits.MaxSteam / 2, limits.MaxReflux / 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := limits.ValidateEnergy(Reboiler{Steam: tc.steam}, Reflux{Ratio: tc.ratio})
			if err != nil {
				t.Fatalf("expected boundary values to pass, got %v", err)
			}
		})
	}
}
