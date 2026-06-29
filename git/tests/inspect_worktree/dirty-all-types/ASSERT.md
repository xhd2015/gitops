## Expected

- `IsClean` is false
- `Added=1`, `Changed=1`, `Renamed=1`, `Deleted=1`
- `Uncommitted=4`

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inspect := resp.Inspect
	if inspect.IsClean {
		t.Fatal("expected dirty worktree")
	}
	if inspect.Added != 1 {
		t.Fatalf("Added = %d, want 1", inspect.Added)
	}
	if inspect.Changed != 1 {
		t.Fatalf("Changed = %d, want 1", inspect.Changed)
	}
	if inspect.Renamed != 1 {
		t.Fatalf("Renamed = %d, want 1", inspect.Renamed)
	}
	if inspect.Deleted != 1 {
		t.Fatalf("Deleted = %d, want 1", inspect.Deleted)
	}
	if inspect.Uncommitted != 4 {
		t.Fatalf("Uncommitted = %d, want 4", inspect.Uncommitted)
	}
}
```