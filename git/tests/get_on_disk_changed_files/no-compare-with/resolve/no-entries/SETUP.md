## Preconditions
- No changed files exist on disk (or all changes are filtered out)

## Steps
1. Create `_keep.go` tracked file to establish shared context

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir
	p := filepath.Join(dir, "_keep.go")
	if err := os.WriteFile(p, []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "_keep.go").Run(); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "commit", "-m", "add _keep.go").Run()
}
```
