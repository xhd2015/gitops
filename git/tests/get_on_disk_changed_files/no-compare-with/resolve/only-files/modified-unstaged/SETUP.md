## Steps
1. Modify `README.md` (already tracked) with new content "modified content"

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	return os.WriteFile(filepath.Join(req.Dir, "README.md"), []byte("modified content"), 0644)
}
```