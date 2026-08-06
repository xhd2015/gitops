## Expected

- `Run` returns a non-nil error (the parse error itself).
- `Error()` is non-empty and mentions `parse` or `diff` (case-insensitive).
- `Unwrap` / cause message equals `req.CauseMsg`.
- `errors.As` succeeds and recovers `Dir` and `Raw` from the typed error.

## Errors

- Typed `DiffCachedParseError` with wrapped cause.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected non-nil error from DiffCachedParseError")
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.ErrMsg == "" {
		t.Fatal("Error() returned empty string")
	}
	lower := strings.ToLower(resp.ErrMsg)
	if !strings.Contains(lower, "parse") && !strings.Contains(lower, "diff") {
		t.Fatalf("Error() = %q, want it to mention parse or diff", resp.ErrMsg)
	}
	if resp.UnwrappedMsg != req.CauseMsg {
		t.Fatalf("Unwrap() = %q, want %q", resp.UnwrappedMsg, req.CauseMsg)
	}
	if !resp.AsOK {
		t.Fatal("errors.As failed to recover *git.DiffCachedParseError")
	}
	if resp.AsDir != req.Dir {
		t.Fatalf("As Dir = %q, want %q", resp.AsDir, req.Dir)
	}
	if resp.AsRaw != req.Raw {
		t.Fatalf("As Raw = %q, want %q", resp.AsRaw, req.Raw)
	}
}
```
