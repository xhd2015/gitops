## Expected
- Result contains `mod.go` (del.go excluded)

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "mod.go" {
		t.Fatalf("expected [mod.go], got: %v", resp.Files)
	}
}
```
