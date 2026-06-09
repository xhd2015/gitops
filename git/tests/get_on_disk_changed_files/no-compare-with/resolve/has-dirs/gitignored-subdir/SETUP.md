## Steps
1. Write `.gitignore` with content `build/` in the repo root
2. Create directories `view/` and `view/build/` under the repo root
3. Write `view/a.go` with content "package view"
4. Write `view/build/output.js` with content "console.log('test')"

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("build/"), 0644); err != nil {
		return err
	}
	vdir := filepath.Join(dir, "view")
	bdir := filepath.Join(vdir, "build")
	if err := os.MkdirAll(bdir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(vdir, "a.go"), []byte("package view"), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(bdir, "output.js"), []byte("console.log('test')"), 0644)
}
```