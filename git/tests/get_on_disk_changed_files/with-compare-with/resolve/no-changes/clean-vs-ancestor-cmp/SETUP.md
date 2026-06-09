## Steps
1. Set req.CompareWith = "HEAD~1"
2. Write `a.go`, add and commit (commit A)
3. Run `git commit --allow-empty -m "empty"` (commit B = HEAD)
4. Working tree is clean

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	req.CompareWith = "HEAD~1"
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "a.go").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "commit A").Run(); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "commit", "--allow-empty", "-m", "empty").Run()
}
```