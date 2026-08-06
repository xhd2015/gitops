## Expected
- Result contains `staged.go`

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "staged.go" {
		t.Fatalf("expected [staged.go], got: %v", resp.Files)
	}
}
```
