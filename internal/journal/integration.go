package journal

import (
	"errors"
	"github.com/wyw14/solo-039/internal/model"
)

func PersistProof(store *Store, proof model.Proof, campaign model.Campaign) error {
	if proof.Digest == "" {
		proof.Digest = model.ProofDigest(proof)
	}
	event := model.NewEvent("proof.durable", campaign, proof.Digest)
	if err := store.Append(event); err != nil {
		return err
	}
	return nil
}

func RequireProof(store *Store, revision string) error {
	event, ok := store.Last("proof.durable")
	if !ok || event.Revision.String() != revision {
		return errors.New("proof is not durable")
	}
	return nil
}
