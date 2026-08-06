## Steps
1. Run `git mv README.md new.go`

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"os/exec"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	return exec.Command("git", "-C", req.Dir, "mv", "README.md", "new.go").Run()
}
```