package interlock

import (
	"github.com/wyw14/solo-039/internal/model"
	"sort"
	"strings"
)

type Explanation struct {
	Allowed bool
	Sources []string
	Reasons []string
	Summary string
}

func Explain(permit Permit) Explanation {
	sources := make([]string, 0, len(permit.Evidence))
	for _, proof := range permit.Evidence {
		sources = append(sources, proof.Source)
	}
	sort.Strings(sources)
	reasons := append([]string(nil), permit.Reasons...)
	sort.Strings(reasons)
	summary := "allowed"
	if !permit.Allowed {
		summary = strings.Join(reasons, "; ")
	}
	return Explanation{Allowed: permit.Allowed, Sources: sources, Reasons: reasons, Summary: summary}
}

func EvidenceSources(proofs []model.Proof) []string {
	sources := make([]string, 0, len(proofs))
	for _, proof := range proofs {
		sources = append(sources, proof.Source)
	}
	sort.Strings(sources)
	return sources
}

func SameSource(a, b model.Proof) bool { return a.Source != "" && a.Source == b.Source }
