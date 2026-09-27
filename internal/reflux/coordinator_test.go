package reflux

import "testing"

func TestCoordinatorBoundsSequentialCapacity(t *testing.T) {
	coordinator := NewCoordinator(1)
	if _, err := coordinator.Acquire("T-101", 0.7, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := coordinator.Acquire("T-102", 0.7, 1); err == nil {
		t.Fatal("second allocation exceeded capacity")
	}
	if coordinator.Used() > 1 {
		t.Fatalf("capacity exceeded: %f", coordinator.Used())
	}
}
