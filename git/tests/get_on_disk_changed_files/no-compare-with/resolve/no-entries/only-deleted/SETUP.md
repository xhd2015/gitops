## Steps
1. Run `git rm README.md` to stage a delete

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"os/exec"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	return exec.Command("git", "-C", req.Dir, "rm", "README.md").Run()
}
```