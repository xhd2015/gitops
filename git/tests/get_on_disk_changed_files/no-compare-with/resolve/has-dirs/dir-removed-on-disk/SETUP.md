## Steps
1. Create `_removed/` directory, then delete it from disk to simulate dir removed between status and expand

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	rdir := filepath.Join(dir, "_removed")
	if err := os.MkdirAll(rdir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(rdir, "f.go"), []byte("package main"), 0644); err != nil {
		return err
	}
	return os.RemoveAll(rdir)
}
```