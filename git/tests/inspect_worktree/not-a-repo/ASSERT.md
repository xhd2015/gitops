## Expected

- `IsRepo` is false
- Branch, commit, and message are empty
- `IsClean` is true; all counts are zero
- No error

## Exit Code

N/A (library call)

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inspect := resp.Inspect
	if inspect.IsRepo {
		t.Fatal("IsRepo = true, want false")
	}
	if inspect.Branch != "" || inspect.CommitShort != "" || inspect.CommitMessage != "" {
		t.Fatalf("expected empty branch/commit, got %+v", inspect)
	}
	if !inspect.IsClean || inspect.Uncommitted != 0 {
		t.Fatalf("expected clean zero counts, got %+v", inspect)
	}
	if inspect.Added != 0 || inspect.Changed != 0 || inspect.Renamed != 0 || inspect.Deleted != 0 {
		t.Fatalf("expected zero type counts, got %+v", inspect)
	}
}
```