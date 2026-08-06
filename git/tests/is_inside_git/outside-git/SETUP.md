## Steps
1. Create a temp directory WITHOUT `git init`
2. Set `req.Dir` to that directory

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    dir, err := os.MkdirTemp("", "notgit")
    if err != nil {
        return err
    }
    t.Cleanup(func() { os.RemoveAll(dir) })
    req.Dir = dir
    return nil
}
```
