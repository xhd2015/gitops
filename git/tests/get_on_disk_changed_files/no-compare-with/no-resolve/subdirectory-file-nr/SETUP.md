## Steps
1. Create directory `sub/pkg/`
2. Write `sub/pkg/foo.go` with content "original"
3. Run `git add -A` and `git commit -m "add subdir file"`
4. Modify `sub/pkg/foo.go` with content "modified"

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	p := filepath.Join(dir, "sub", "pkg", "foo.go")
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(p, []byte("original"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "-A").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "add subdir file").Run(); err != nil {
		return err
	}
	return os.WriteFile(p, []byte("modified"), 0644)
}
```