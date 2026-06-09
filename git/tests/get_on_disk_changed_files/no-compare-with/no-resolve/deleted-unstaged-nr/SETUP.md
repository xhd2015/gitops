## Steps
1. Remove `README.md` from disk (e.g. `os.Remove("README.md")`) without staging

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	return os.Remove(filepath.Join(req.Dir, "README.md"))
}
```