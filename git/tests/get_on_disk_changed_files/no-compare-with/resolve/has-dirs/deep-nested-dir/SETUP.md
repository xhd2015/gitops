## Steps
1. Write `.gitignore` with content `*.log` in the repo root
2. Create directory `a/b/c/` under the repo root
3. Write `a/b/c/d.go` with content "package d"
4. Write `a/b/ignored.log` with content "log content"

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
	cdir := filepath.Join(dir, "a", "b", "c")
	if err := os.MkdirAll(cdir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(cdir, "d.go"), []byte("package d"), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "a", "b", "ignored.log"), []byte("log content"), 0644)
}
```