## Expected
- An error is returned because the directory is not a git repository

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for non-git directory, got nil")
	}
}
```
