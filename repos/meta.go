package repos

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type cloneMeta struct {
	CloneURL string `json:"cloneURL"`
	Depth    int    `json:"depth"`
}

func metaPath(cacheDir string) string {
	return cacheDir + ".meta"
}

func writeCloneMeta(cacheDir, cloneURL string, opts BareCacheOptions) error {
	meta := cloneMeta{
		CloneURL: cloneURL,
		Depth:    opts.Depth,
	}
	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(cacheDir), 0o755); err != nil {
		return err
	}
	return os.WriteFile(metaPath(cacheDir), data, 0o644)
}

func readCloneMeta(cacheDir string) (*cloneMeta, error) {
	data, err := os.ReadFile(metaPath(cacheDir))
	if err != nil {
		return nil, err
	}
	var meta cloneMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	if meta.CloneURL == "" {
		return nil, fmt.Errorf("clone meta missing URL for %s", cacheDir)
	}
	return &meta, nil
}

func removeCloneMeta(cacheDir string) {
	_ = os.Remove(metaPath(cacheDir))
}