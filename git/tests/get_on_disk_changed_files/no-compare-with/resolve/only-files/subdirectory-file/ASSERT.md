## Expected
- Result contains `sub/pkg/foo.go`

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "sub/pkg/foo.go" {
		t.Fatalf("expected [sub/pkg/foo.go], got: %v", resp.Files)
	}
}
```
