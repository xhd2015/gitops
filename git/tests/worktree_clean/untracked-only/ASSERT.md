## Expected

- `DiffClean` is true (untracked ignored by diff-vs-HEAD)
- `PorcelainClean` is false (untracked counts as dirty)
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
		t.Fatal("IsDiffClean = false, want true for untracked-only")
	}
	if resp.PorcelainClean {
		t.Fatal("IsPorcelainClean = true, want false for untracked-only")
	}
}
```
