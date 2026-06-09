## Steps
1. Set req.CompareWith = "HEAD"
2. Make at least one additional commit
3. Write `a.go` and `b.go`, add and commit
4. Modify both `a.go` and `b.go`

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
	for _, name := range []string{"a.go", "b.go"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("original"), 0644); err != nil {
			return err
		}
	}
	if err := exec.Command("git", "-C", dir, "add", "a.go", "b.go").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "add files").Run(); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("modified"), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "b.go"), []byte("modified"), 0644)
}
```