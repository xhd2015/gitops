## Expected

- All `ChangeCounts` fields are zero
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
	if c.Added != 0 || c.Changed != 0 || c.Renamed != 0 || c.Deleted != 0 {
		t.Fatalf("expected zero ChangeCounts, got %+v", c)
	}
}
```
