## Expected

- After `git restore --staged -- .`, `DiffCached` returns `(nil, nil)`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error after unstage-all: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Diff != nil {
		t.Fatalf("expected Diff == nil after unstage-all, got %#v", resp.Diff)
	}
}
```
