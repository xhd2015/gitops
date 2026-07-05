# Scenario

Repo-local core.hooksPath still runs under isolated git.

## Steps

1. Create repo with failing pre-commit hook via core.hooksPath.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	repo := filepath.Join(t.TempDir(), "repo")
	hooksDir := filepath.Join(t.TempDir(), "hooks")
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(hooksDir, "pre-commit"), []byte("#!/bin/sh\nexit 2\n"), 0o755); err != nil {
		return err
	}
	req.RepoDir = repo
	req.HooksDir = hooksDir
	return nil
}
```