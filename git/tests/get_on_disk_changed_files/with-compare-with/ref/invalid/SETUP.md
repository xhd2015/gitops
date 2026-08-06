## Steps
1. Set req.CompareWith = "nonexistent123"

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.CompareWith = "nonexistent123"
	return nil
}
```