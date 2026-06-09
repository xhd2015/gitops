## Steps
1. Get HEAD~1 commit hash: `git rev-parse HEAD~1`
2. Set req.CompareWith to that commit hash
3. Write `a.go` with content "package main", add and commit
4. Working tree is clean

```go
import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	out, err := exec.Command("git", "-C", dir, "rev-parse", "HEAD~1").Output()
	if err != nil {
		return fmt.Errorf("rev-parse HEAD~1 failed: %w", err)
	}
	req.CompareWith = string(bytes.TrimSpace(out))
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