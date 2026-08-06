## Steps
- Create a temporary directory that is NOT a git repository
- Set req.Dir to this non-git directory
- Set RefA and RefB to arbitrary values — neither will resolve since there is no git repo

```go
import (
	"os"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir, err := os.MkdirTemp("", "gitops-compare-nogit-*")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	req.Dir = dir
	req.RefA = "main"
	req.RefB = "main"
	return nil
}
```
