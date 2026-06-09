## Steps
1. Write `file` as a regular file with content "regular"
2. Run `git add file` and `git commit -m "add file"`
3. Create a symlink: `ln -sf /tmp file` (or equivalent to trigger type change T)

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	p := filepath.Join(dir, "file")
	if err := os.WriteFile(p, []byte("regular"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "file").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "add file").Run(); err != nil {
		return err
	}
	if err := os.Remove(p); err != nil {
		return err
	}
	return os.Symlink("/tmp", p)
}
```