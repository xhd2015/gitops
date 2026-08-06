## Steps
1. Modify `README.md` with new content "staged content"
2. Run `git add README.md` to stage the change

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("staged content"), 0644); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "add", "README.md").Run()
}
```