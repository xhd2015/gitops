## Steps
1. Set req.CompareWith = "HEAD"
2. Make at least one additional commit
3. Write `mod.go` with content "original", add and commit
4. Modify `mod.go` to "modified"

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