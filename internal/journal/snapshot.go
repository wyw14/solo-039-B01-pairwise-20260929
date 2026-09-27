package journal

import (
	"encoding/json"
	"github.com/wyw14/solo-039/internal/model"
	"os"
	"path/filepath"
)

type SnapshotFile struct{ path string }

func NewSnapshotFile(path string) *SnapshotFile { return &SnapshotFile{path: path} }

func (f *SnapshotFile) Save(snapshot model.Snapshot) error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0755); err != nil {
		return err
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	temporary := f.path + ".tmp"
	if err := os.WriteFile(temporary, data, 0644); err != nil {
		return err
	}
	file, err := os.OpenFile(temporary, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(temporary, f.path)
}

func (f *SnapshotFile) Load() (model.Snapshot, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return model.Snapshot{}, err
	}
	var snapshot model.Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return model.Snapshot{}, err
	}
	return snapshot, nil
}
