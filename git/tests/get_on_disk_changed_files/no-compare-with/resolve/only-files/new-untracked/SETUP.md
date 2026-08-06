## Steps
1. Write `new.go` with content "package main"

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	return os.WriteFile(filepath.Join(req.Dir, "new.go"), []byte("package main"), 0644)
}
```