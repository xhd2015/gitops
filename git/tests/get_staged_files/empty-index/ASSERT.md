## Expected

- `GetStagedFiles` returns an empty slice.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 0 {
		t.Fatalf("expected 0 staged files, got %d: %v", len(resp.Files), resp.Files)
	}
}
```
