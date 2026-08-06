## Expected
- Result is nil (no file changes between HEAD~1 and HEAD, both clean commits)

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
		t.Fatalf("expected nil (no diff between clean commits), got: %v", resp.Files)
	}
}
```
