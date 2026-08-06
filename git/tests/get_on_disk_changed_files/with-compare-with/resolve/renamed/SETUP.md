## Preconditions
- There are renamed files (RenamedFrom != "")

## Steps
1. Create `_rename_base.go` tracked file to establish shared context

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir
	p := filepath.Join(dir, "_rename_base.go")
	if err := os.WriteFile(p, []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "_rename_base.go").Run(); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "commit", "-m", "add _rename_base.go").Run()
}
```
