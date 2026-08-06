## Expected
- `resp.Inside` is `true`
- No error

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !resp.Inside {
        t.Fatal("expected Inside=true, got false")
    }
}
```
