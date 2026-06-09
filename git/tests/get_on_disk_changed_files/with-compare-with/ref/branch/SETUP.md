## Steps
1. Set req.CompareWith = "master"
2. Make at least one additional commit
3. Working tree is clean

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	req.CompareWith = "master"
	if err := os.WriteFile(filepath.Join(dir, "_dummy.go"), []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "_dummy.go").Run(); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "commit", "-m", "dummy commit").Run()
}
```