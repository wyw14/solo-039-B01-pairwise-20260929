package model

import "errors"

var transitionTable = map[Phase]map[Phase]bool{
	Preparing:     {Warming: true, Stopped: true},
	Warming:       {Stabilizing: true, Shutting: true, Stopped: true},
	Stabilizing:   {OnSpec: true, Transitioning: true, Shutting: true},
	OnSpec:        {Transitioning: true, Shutting: true},
	Transitioning: {Warming: true, Stabilizing: true, Shutting: true},
	Shutting:      {Stopped: true},
	Stopped:       {Preparing: true},
}

func CanTransition(from, to Phase) bool { return transitionTable[from][to] }

func TransitionError(from, to Phase) error {
	if CanTransition(from, to) {
		return nil
	}
	return errors.New("invalid campaign transition")
}

func Phases() []Phase {
	return []Phase{Preparing, Warming, Stabilizing, OnSpec, Transitioning, Shutting, Stopped}
}
