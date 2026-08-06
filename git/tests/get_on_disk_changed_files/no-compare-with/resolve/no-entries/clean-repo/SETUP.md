## Steps
1. No additional steps. The repo is clean after initial commit.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Log("repo is clean after initial commit — no additional setup needed")
	return nil
}
```