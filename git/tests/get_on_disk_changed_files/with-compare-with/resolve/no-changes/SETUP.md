## Preconditions
- Working tree is clean relative to compare ref

## Steps
1. Create `_nochg.go` tracked file to establish shared context

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	p := filepath.Join(dir, "_nochg.go")
	if err := os.WriteFile(p, []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "_nochg.go").Run(); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "commit", "-m", "add _nochg.go").Run()
}
```
