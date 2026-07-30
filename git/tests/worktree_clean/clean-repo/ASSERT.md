## Expected

- `DiffClean` is true
- `PorcelainClean` is true
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
	if !resp.DiffClean {
		t.Fatal("IsDiffClean = false, want true on clean repo")
	}
	if !resp.PorcelainClean {
		t.Fatal("IsPorcelainClean = false, want true on clean repo")
	}
}
```
