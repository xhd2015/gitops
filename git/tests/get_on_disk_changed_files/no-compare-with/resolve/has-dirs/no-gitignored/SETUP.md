## Steps
1. Create directory `view/` under the repo root
2. Write `view/a.go` with content "package view"
3. Write `view/b.go` with content "package view"

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	vdir := filepath.Join(dir, "view")
	if err := os.MkdirAll(vdir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(vdir, "a.go"), []byte("package view"), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(vdir, "b.go"), []byte("package view"), 0644)
}
```