## Steps
1. Set req.CompareWith = "HEAD"
2. Make at least one additional commit so HEAD != initial commit
3. Write `new.go` with content "package main" (untracked)

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
	return os.WriteFile(filepath.Join(dir, "new.go"), []byte("package main"), 0644)
}
```