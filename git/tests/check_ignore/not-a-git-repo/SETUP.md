## Steps
1. Set `req.Dir` to a directory that is NOT a git repository
2. Set `req.Path = "any.file"`

```go
import (
	"os"
)

func Setup(t *testing.T, req *Request) error {
	dir, err := os.MkdirTemp("", "not-a-git-repo")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	req.Dir = dir
	req.Path = "any.file"
	return nil
}
```
