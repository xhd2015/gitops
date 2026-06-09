## Steps
1. Set req.CompareWith = "HEAD"
2. Make at least one additional commit
3. Create directory `view/`
4. Write `view/a.go` with content "package view"

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
	vdir := filepath.Join(dir, "view")
	if err := os.MkdirAll(vdir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(vdir, "a.go"), []byte("package view"), 0644)
}
```