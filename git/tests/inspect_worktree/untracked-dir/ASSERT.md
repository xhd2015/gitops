## Expected

- `Added=4` (1 file + 3 under nested/)
- `Uncommitted=2` (porcelain lines unchanged)

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inspect := resp.Inspect
	if inspect.Added != 4 {
		t.Fatalf("Added = %d, want 4", inspect.Added)
	}
	if inspect.Uncommitted != 2 {
		t.Fatalf("Uncommitted = %d, want 2", inspect.Uncommitted)
	}
}
```