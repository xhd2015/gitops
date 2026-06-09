## Steps
1. Set req.CompareWith = "HEAD~1"
2. Write `a.go` with content "package main", add and commit
3. Working tree is clean

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	req.CompareWith = "HEAD~1"
	p := filepath.Join(dir, "a.go")
	if err := os.WriteFile(p, []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "a.go").Run(); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "commit", "-m", "add a.go").Run()
}
```