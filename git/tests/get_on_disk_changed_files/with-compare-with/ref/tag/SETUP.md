## Steps
1. Create tag `v1.0` at HEAD of base commit: `git tag v1.0`
2. Set req.CompareWith = "v1.0"
3. Write `a.go` with content "package main", add and commit
4. Working tree is clean

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	if err := exec.Command("git", "-C", dir, "tag", "v1.0").Run(); err != nil {
		return err
	}
	req.CompareWith = "v1.0"
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