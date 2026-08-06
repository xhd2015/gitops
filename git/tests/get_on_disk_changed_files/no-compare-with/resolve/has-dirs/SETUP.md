## Preconditions
- There is at least one untracked directory on disk
- git status --porcelain reports directory entries (e.g. `?? view/`)

## Steps
1. Create `_base/` directory tree to establish a shared untracked directory context for all children

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := filepath.Join(req.Dir, "_base")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "base.go"), []byte("package base"), 0644)
}
```
