## Expected

- No error
- `MainRepo` equals the main checkout path

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !samePath(t, resp.MainRepo, req.MainPath) {
		t.Fatalf("ResolveMainRepo = %q, want main %q", resp.MainRepo, req.MainPath)
	}
}
```
