package repos

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ListCachedRepoKeys returns relative path keys for bare repos under ~/.repos/cache.
func ListCachedRepoKeys() ([]string, error) {
	root, err := ReposRoot()
	if err != nil {
		return nil, err
	}
	cacheRoot := filepath.Join(root, "cache")
	return listRepoKeys(cacheRoot, cacheRoot)
}

func listRepoKeys(cacheRoot, dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var keys []string
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if isBare, err := isBareRepository(path); err == nil && isBare {
				rel, err := filepath.Rel(cacheRoot, path)
				if err != nil {
					return nil, err
				}
				keys = append(keys, filepath.ToSlash(rel))
				continue
			}
			nested, err := listRepoKeys(cacheRoot, path)
			if err != nil {
				return nil, err
			}
			keys = append(keys, nested...)
		}
	}
	return keys, nil
}

// FormatCachedRepoKey returns a human-readable cache key (slashes as underscores).
func FormatCachedRepoKey(relPath string) string {
	return strings.ReplaceAll(relPath, "/", "_")
}

// CacheDirFromKey resolves a formatted cache key back to a bare cache directory.
func CacheDirFromKey(key string) (string, error) {
	root, err := ReposRoot()
	if err != nil {
		return "", err
	}
	rel := strings.ReplaceAll(key, "_", "/")
	return filepath.Join(append([]string{root, "cache"}, strings.Split(rel, "/")...)...), nil
}

// ValidateCachedRepoKey reports whether key maps to an existing bare cache.
func ValidateCachedRepoKey(key string) (string, error) {
	dir, err := CacheDirFromKey(key)
	if err != nil {
		return "", err
	}
	isBare, err := isBareRepository(dir)
	if err != nil {
		return "", err
	}
	if !isBare {
		return "", fmt.Errorf("not a bare repository: %s", dir)
	}
	return dir, nil
}