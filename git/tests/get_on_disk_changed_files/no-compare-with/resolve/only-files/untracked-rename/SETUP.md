## Steps
1. Run `os.Rename(README.md, RENAMED.go)` (not git mv — plain filesystem rename)

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	return os.Rename(filepath.Join(dir, "README.md"), filepath.Join(dir, "RENAMED.go"))
}
```