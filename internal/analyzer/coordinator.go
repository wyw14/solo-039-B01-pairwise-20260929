package analyzer

import (
	"github.com/google/uuid"
	"github.com/wyw14/solo-039/internal/model"
)

type Coordinator struct{ State *State }

func NewCoordinator() *Coordinator { return &Coordinator{State: New()} }

func (c *Coordinator) Begin(revision uuid.UUID) { c.State.BindFeed(revision) }

func (c *Coordinator) Record(revision, calibration uuid.UUID, value float64, at int64) bool {
	return c.State.Add(model.Sample{ID: uuid.New(), FeedRevision: revision, Calibration: calibration, Value: value, Component: "top", At: at, Valid: true})
}

func (c *Coordinator) CurrentProof(revision uuid.UUID) (model.Proof, bool) {
	return c.State.Proof(revision)
}
