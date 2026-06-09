## Preconditions
- All changed paths are individual files (no directory entries in git-status)

## Steps
1. Create `_onlyfiles.go` tracked file to establish shared context

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	p := filepath.Join(dir, "_onlyfiles.go")
	if err := os.WriteFile(p, []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "_onlyfiles.go").Run(); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "commit", "-m", "add _onlyfiles.go").Run()
}
```
