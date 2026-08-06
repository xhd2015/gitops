## Steps
1. Set req.Dir to a path that does not exist on disk

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Dir = "/nonexistent/path/that/does/not/exist"
	return nil
}
```
