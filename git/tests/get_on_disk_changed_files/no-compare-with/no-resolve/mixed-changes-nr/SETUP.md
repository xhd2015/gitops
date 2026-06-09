## Steps
1. Write `mod.go` with content "original", add and commit
2. Write `del.go` with content "original", add and commit
3. Write `ren.go` with content "original", add and commit
4. Modify `mod.go` to "modified"
5. Write `added.go` with content "new staged" and `git add added.go`
6. Write `untracked.go` with content "untracked"
7. Delete `del.go` from disk
8. Run `git mv ren.go ren_new.go`

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	for _, name := range []string{"mod.go", "del.go", "ren.go"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("original"), 0644); err != nil {
			return err
		}
	}
	if err := exec.Command("git", "-C", dir, "add", "mod.go", "del.go", "ren.go").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "add files").Run(); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "mod.go"), []byte("modified"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "added.go"), []byte("new staged"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "added.go").Run(); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "untracked.go"), []byte("untracked"), 0644); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(dir, "del.go")); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "mv", "ren.go", "ren_new.go").Run()
}
```