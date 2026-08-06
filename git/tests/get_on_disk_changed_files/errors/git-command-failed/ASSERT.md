## Expected
- An error is returned because git commands fail against the corrupted repository

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for corrupted git repo, got nil")
	}
}
```
