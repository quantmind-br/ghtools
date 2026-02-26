package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/diogo/ghtools/internal/types"
)

func cachePath() string {
	return filepath.Join(os.TempDir(), "ghtools_repos.json")
}

func IsValid(ttl int) bool {
	info, err := os.Stat(cachePath())
	if err != nil {
		return false
	}
	age := time.Since(info.ModTime())
	return age.Seconds() < float64(ttl)
}

func Read() ([]types.Repo, error) {
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return nil, err
	}
	var repos []types.Repo
	if err := json.Unmarshal(data, &repos); err != nil {
		return nil, err
	}
	return repos, nil
}

func Write(repos []types.Repo) error {
	data, err := json.Marshal(repos)
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath(), data, 0644)
}

func Clear() error {
	return os.Remove(cachePath())
}

func ModTime() time.Time {
	info, err := os.Stat(cachePath())
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
