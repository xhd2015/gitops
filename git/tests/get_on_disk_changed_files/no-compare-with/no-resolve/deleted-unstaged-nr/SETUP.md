## Steps
1. Remove `README.md` from disk (e.g. `os.Remove("README.md")`) without staging

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	return os.Remove(filepath.Join(req.Dir, "README.md"))
}
```