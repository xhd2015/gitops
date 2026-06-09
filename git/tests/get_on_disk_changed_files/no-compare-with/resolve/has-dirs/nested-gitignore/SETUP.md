## Steps
1. Create directory `view/sub/` under the repo root
2. Write `view/sub/.gitignore` with content `*.tmp`
3. Write `view/sub/keep.go` with content "package sub"
4. Write `view/sub/a.tmp` with content "temporary data"

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	sdir := filepath.Join(dir, "view", "sub")
	if err := os.MkdirAll(sdir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(sdir, ".gitignore"), []byte("*.tmp"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(sdir, "keep.go"), []byte("package sub"), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(sdir, "a.tmp"), []byte("temporary data"), 0644)
}
```