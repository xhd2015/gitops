## Expected
- Result contains `README.md`

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "README.md" {
		t.Fatalf("expected [README.md], got: %v", resp.Files)
	}
}
```
