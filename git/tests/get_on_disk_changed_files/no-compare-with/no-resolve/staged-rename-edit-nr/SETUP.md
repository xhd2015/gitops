## Steps
1. Run `git mv README.md new.go`
2. Modify `new.go` with content "modified after rename"

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	if err := exec.Command("git", "-C", dir, "mv", "README.md", "new.go").Run(); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "new.go"), []byte("modified after rename"), 0644)
}
```