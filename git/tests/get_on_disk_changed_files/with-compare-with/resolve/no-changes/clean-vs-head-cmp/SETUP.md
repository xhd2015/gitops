## Steps
1. Set req.CompareWith = "HEAD"
2. Make at least one additional commit (so HEAD is not initial)
3. Working tree is clean (no modifications)

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
	return exec.Command("git", "-C", dir, "commit", "-m", "dummy commit").Run()
}
```