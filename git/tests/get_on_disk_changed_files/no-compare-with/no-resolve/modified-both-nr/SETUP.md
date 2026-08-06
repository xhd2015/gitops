## Steps
1. Modify `README.md` with content "v1"
2. Run `git add README.md`
3. Modify `README.md` again with content "v2"

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir
	p := filepath.Join(dir, "README.md")
	if err := os.WriteFile(p, []byte("v1"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", "README.md").Run(); err != nil {
		return err
	}
	return os.WriteFile(p, []byte("v2"), 0644)
}
```