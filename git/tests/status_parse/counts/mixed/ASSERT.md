## Expected

- `Modified=2`, `Untracked=1`
- Other count fields zero
- No error

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c := resp.Counts
	if c.Modified != 2 {
		t.Fatalf("Modified = %d, want 2", c.Modified)
	}
	if c.Untracked != 1 {
		t.Fatalf("Untracked = %d, want 1", c.Untracked)
	}
	if c.Added != 0 || c.Deleted != 0 || c.Renamed != 0 || c.Copied != 0 || c.Unmerged != 0 {
		t.Fatalf("unexpected non-zero fields: %+v", c)
	}
}
```
