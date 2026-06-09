## Steps
1. Set req.CompareWith = "HEAD"
2. Make at least one additional commit
3. Write `mod.go`, `del.go`, and `keep.go`, add and commit
4. Modify `mod.go` to "modified"
5. Delete `del.go` from disk
6. Write `new.go` with content "package main" (untracked)

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
	for _, name := range []string{"mod.go", "del.go", "keep.go"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("original"), 0644); err != nil {
			return err
		}
	}
	if err := exec.Command("git", "-C", dir, "add", "mod.go", "del.go", "keep.go").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "add files").Run(); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "mod.go"), []byte("modified"), 0644); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(dir, "del.go")); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "new.go"), []byte("package main"), 0644)
}
```