## Steps
1. Write `staged.go` with content "package main"
2. Run `git add staged.go`

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir
	if err := os.WriteFile(filepath.Join(dir, "staged.go"), []byte("package main"), 0644); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "add", "staged.go").Run()
}
```