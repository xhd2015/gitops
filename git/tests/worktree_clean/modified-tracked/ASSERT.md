## Expected

- `DiffClean` is false
- `PorcelainClean` is false
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
	if resp.DiffClean {
		t.Fatal("IsDiffClean = true, want false for modified tracked file")
	}
	if resp.PorcelainClean {
		t.Fatal("IsPorcelainClean = true, want false for modified tracked file")
	}
}
```
