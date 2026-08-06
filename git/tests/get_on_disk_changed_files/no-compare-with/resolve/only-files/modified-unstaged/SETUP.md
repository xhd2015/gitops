## Steps
1. Modify `README.md` (already tracked) with new content "modified content"

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	return os.WriteFile(filepath.Join(req.Dir, "README.md"), []byte("modified content"), 0644)
}
```