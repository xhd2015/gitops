## Steps
1. Create directory `view/` under the repo root
2. Write `view/a.go` with content "package view"

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir
	if err := os.MkdirAll(filepath.Join(dir, "view"), 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "view", "a.go"), []byte("package view"), 0644)
}
```