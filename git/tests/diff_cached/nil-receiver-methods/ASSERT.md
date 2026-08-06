## Expected

- Run succeeds with `resp.Diff == nil` (nil-methods mode).
- On a nil `*model.CachedDiff` receiver:
  - `FileCount() == 0`
  - `Unified() == ""`
  - `UnifiedTruncated(24) == ""`

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/gitops/model"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error for nil-methods mode: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Diff != nil {
		t.Fatalf("expected Diff == nil for nil-methods mode, got %#v", resp.Diff)
	}

	var d *model.CachedDiff // nil receiver under test
	if got := d.FileCount(); got != 0 {
		t.Fatalf("nil FileCount() = %d, want 0", got)
	}
	if got := d.Unified(); got != "" {
		t.Fatalf("nil Unified() = %q, want empty", got)
	}
	if got := d.UnifiedTruncated(24); got != "" {
		t.Fatalf("nil UnifiedTruncated(24) = %q, want empty", got)
	}
}
```
