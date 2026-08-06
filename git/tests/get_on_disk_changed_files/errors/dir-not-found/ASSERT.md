## Expected
- An error is returned because the directory does not exist

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for nonexistent directory, got nil")
	}
}
```
