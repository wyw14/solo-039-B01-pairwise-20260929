package analyzer

import (
	"errors"
	"github.com/wyw14/solo-039/internal/model"
	"sort"
	"sync"
)

type EvidenceCatalog struct {
	mu         sync.Mutex
	byRevision map[string][]model.Proof
}

func NewEvidenceCatalog() *EvidenceCatalog {
	return &EvidenceCatalog{byRevision: make(map[string][]model.Proof)}
}

func (c *EvidenceCatalog) Add(proof model.Proof) error {
	if proof.Revision.String() == "00000000-0000-0000-0000-000000000000" {
		return errors.New("proof revision is missing")
	}
	if proof.Source == "" || proof.Digest == "" {
		return errors.New("proof content is incomplete")
	}
	c.mu.Lock()
	c.byRevision[proof.Revision.String()] = append(c.byRevision[proof.Revision.String()], proof)
	c.mu.Unlock()
	return nil
}

func (c *EvidenceCatalog) Independent(revision string) []model.Proof {
	c.mu.Lock()
	defer c.mu.Unlock()
	values := c.byRevision[revision]
	seen := make(map[string]struct{})
	out := make([]model.Proof, 0, len(values))
	for _, proof := range values {
		if !proof.Independent {
			continue
		}
		if _, ok := seen[proof.Source]; ok {
			continue
		}
		seen[proof.Source] = struct{}{}
		out = append(out, proof)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].At < out[j].At })
	return out
}

func (c *EvidenceCatalog) ClearExcept(revision string) {
	c.mu.Lock()
	values := append([]model.Proof(nil), c.byRevision[revision]...)
	c.byRevision = map[string][]model.Proof{revision: values}
	c.mu.Unlock()
}

func (c *EvidenceCatalog) Count(revision string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.byRevision[revision])
}
