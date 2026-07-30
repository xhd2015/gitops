## Expected

- `Added=1` (`??`)
- `Changed=1` (` M`)
- `Renamed=1` (`R `)
- `Deleted=1` (` D`)
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
	c := resp.ChangeCounts
	if c.Added != 1 {
		t.Fatalf("Added = %d, want 1", c.Added)
	}
	if c.Changed != 1 {
		t.Fatalf("Changed = %d, want 1", c.Changed)
	}
	if c.Renamed != 1 {
		t.Fatalf("Renamed = %d, want 1", c.Renamed)
	}
	if c.Deleted != 1 {
		t.Fatalf("Deleted = %d, want 1", c.Deleted)
	}
}
```
