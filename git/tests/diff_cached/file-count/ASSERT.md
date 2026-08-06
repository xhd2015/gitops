## Expected

- `DiffCached` returns a non-nil `*model.CachedDiff` with no error.
- `resp.Diff.FileCount() == 2` for the two staged files.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error for two staged files: %v", err)
	}
	if resp == nil || resp.Diff == nil {
		t.Fatal("expected non-nil Diff for two staged files")
	}
	got := resp.Diff.FileCount()
	if got != 2 {
		t.Fatalf("FileCount() = %d, want 2 (Files=%#v)", got, resp.Diff.Files)
	}
}
```
