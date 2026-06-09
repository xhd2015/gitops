## Steps
1. Set req.CompareWith = "HEAD"
2. Make at least one additional commit
3. Write `old.go` with content "package main", add and commit
4. Run `git mv old.go new.go`

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
	p := filepath.Join(dir, "old.go")
	if err := os.WriteFile(p, []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "old.go").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "add old.go").Run(); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "mv", "old.go", "new.go").Run()
}
```