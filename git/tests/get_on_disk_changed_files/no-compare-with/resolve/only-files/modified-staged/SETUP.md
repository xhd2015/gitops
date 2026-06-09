## Steps
1. Modify `README.md` with new content "staged content"
2. Run `git add README.md` to stage the change

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("staged content"), 0644); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "add", "README.md").Run()
}
```