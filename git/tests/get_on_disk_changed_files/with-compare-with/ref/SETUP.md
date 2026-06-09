## Preconditions
- Tests covering different CompareWith ref types (HEAD, HEAD~1, branch, tag, commit hash, invalid)

## Steps
1. Create `base.go` and commit, to ensure HEAD has a parent commit for ancestor refs

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	if err := os.WriteFile(filepath.Join(dir, "base.go"), []byte("package main"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "base.go").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "base commit").Run(); err != nil {
		return err
	}
	return nil
}
```
