## Steps
1. Write `.gitignore` with content `*.log` in the repo root
2. Create directory `view/` under the repo root
3. Write `view/a.go` with content "package view"
4. Write `view/debug.log` with content "log data"

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.log"), 0644); err != nil {
		return err
	}
	vdir := filepath.Join(dir, "view")
	if err := os.MkdirAll(vdir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(vdir, "a.go"), []byte("package view"), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(vdir, "debug.log"), []byte("log data"), 0644)
}
```