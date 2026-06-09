## Preconditions
- There are newly staged files (IsNew)

## Steps
1. Create `_newbase.go` tracked file to establish shared context

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	p := filepath.Join(dir, "_newbase.go")
	if err := os.WriteFile(p, []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "_newbase.go").Run(); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "commit", "-m", "add _newbase.go").Run()
}
```
