## Expected

- One entry
- Path `/path/to/detached`
- Branch is empty string
- HEAD is set

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Entries) != 1 {
		t.Fatalf("len = %d, want 1", len(resp.Entries))
	}
	e := resp.Entries[0]
	if e.Path != "/path/to/detached" {
		t.Fatalf("Path = %q", e.Path)
	}
	if e.Branch != "" {
		t.Fatalf("Branch = %q, want empty for detached", e.Branch)
	}
	if e.HEAD != "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Fatalf("HEAD = %q", e.HEAD)
	}
}
```
