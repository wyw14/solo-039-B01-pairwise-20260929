package journal

import (
	"github.com/wyw14/solo-039/internal/model"
	"testing"
)

func TestStoreFailureDoesNotPublishEvent(t *testing.T) {
	store := New()
	store.FailNext()
	if err := store.Append(model.NewEvent("proof", model.NewCampaign("g"), "value")); err == nil {
		t.Fatal("expected write failure")
	}
	if len(store.Events()) != 0 {
		t.Fatal("failed event was published")
	}
}
