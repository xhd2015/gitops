## Steps
- Create a temporary directory
- Set req.Dir to a non-existent subdirectory within it
- Set RefA and RefB to arbitrary values

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir, err := os.MkdirTemp("", "gitops-compare-nodir-*")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	req.Dir = filepath.Join(dir, "nonexistent")
	req.RefA = "main"
	req.RefB = "main"
	return nil
}
```
