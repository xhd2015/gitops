## Steps
1. Set req.CompareWith = "HEAD"
2. Make at least one additional commit
3. Write `staged.go` with content "package main"
4. Run `git add staged.go`

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	req.CompareWith = "HEAD"
	if err := os.WriteFile(filepath.Join(dir, "_dummy.go"), []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "_dummy.go").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "dummy commit").Run(); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "staged.go"), []byte("package main"), 0644); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "add", "staged.go").Run()
}
```