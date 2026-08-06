## Steps
1. Set req.CompareWith = "HEAD"
2. Write `mod.go` with content "original", add and commit
3. Modify `mod.go` to "modified"

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir
	req.CompareWith = "HEAD"
	p := filepath.Join(dir, "mod.go")
	if err := os.WriteFile(p, []byte("original"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "mod.go").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "add mod.go").Run(); err != nil {
		return err
	}
	return os.WriteFile(p, []byte("modified"), 0644)
}
```