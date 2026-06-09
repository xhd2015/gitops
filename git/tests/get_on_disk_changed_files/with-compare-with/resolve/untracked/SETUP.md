## Preconditions
- There are untracked files in the working tree

## Steps
1. Create `_utbase.go` tracked file and `_utnew.go` untracked file to establish shared untracked context

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	p := filepath.Join(dir, "_utbase.go")
	if err := os.WriteFile(p, []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "_utbase.go").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "add _utbase.go").Run(); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "_utnew.go"), []byte("package main"), 0644)
}
```
