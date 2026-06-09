## Steps
1. Write `new.go` with content "package main"

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	return os.WriteFile(filepath.Join(req.Dir, "new.go"), []byte("package main"), 0644)
}
```