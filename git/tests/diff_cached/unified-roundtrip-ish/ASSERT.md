## Expected

- `DiffCached` returns a non-nil `*model.CachedDiff` with no error.
- `Unified()` is non-empty.
- `Unified()` contains `diff --git`.
- `Unified()` contains the staged path `note.txt`.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error for staged note.txt: %v", err)
	}
	if resp == nil || resp.Diff == nil {
		t.Fatal("expected non-nil Diff for staged note.txt")
	}
	u := resp.Diff.Unified()
	if u == "" {
		t.Fatal("Unified() returned empty string")
	}
	if !strings.Contains(u, "diff --git") {
		t.Fatalf("Unified() missing %q; got:\n%s", "diff --git", u)
	}
	if !strings.Contains(u, "note.txt") {
		t.Fatalf("Unified() missing path note.txt; got:\n%s", u)
	}
}
```
