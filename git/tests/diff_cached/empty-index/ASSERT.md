## Expected

- `DiffCached` returns `(nil, nil)` — nil `*model.CachedDiff`, no error.
- An empty non-nil struct is not accepted for this contract.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error for empty index: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Diff != nil {
		t.Fatalf("expected Diff == nil for empty staged patch, got %#v", resp.Diff)
	}
}
```
