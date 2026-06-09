## Steps
1. Write `.gitignore` with content `*.o` in the repo root
2. Create directory `build/` under the repo root
3. Write `build/a.o` with content "object"
4. Write `build/b.o` with content "object"

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.o"), 0644); err != nil {
		return err
	}
	bdir := filepath.Join(dir, "build")
	if err := os.MkdirAll(bdir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(bdir, "a.o"), []byte("object"), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(bdir, "b.o"), []byte("object"), 0644)
}
```