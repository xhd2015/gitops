## Expected

- Result Diff is nil (or Response carries nil Diff).
- Error is non-nil and recovers as `*git.DiffCachedParseError`.
- `DiffCachedParseError.Raw` equals the full malformed input (`req.Raw`).

## Errors

- Typed `*DiffCachedParseError` with Raw preserved for diagnostics.

```go
import (
	"errors"
	"testing"

	"github.com/xhd2015/gitops/git"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected parse error for malformed raw, got nil")
	}
	if resp != nil && resp.Diff != nil {
		t.Fatalf("expected nil Diff on parse failure, got %#v", resp.Diff)
	}
	var perr *git.DiffCachedParseError
	if !errors.As(err, &perr) {
		t.Fatalf("expected *git.DiffCachedParseError, got %T: %v", err, err)
	}
	if perr.Raw != req.Raw {
		t.Fatalf("DiffCachedParseError.Raw = %q, want full input %q", perr.Raw, req.Raw)
	}
}
```
