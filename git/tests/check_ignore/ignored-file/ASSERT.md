## Expected
- `CheckIgnore` returns `true` (`app.log` matches `*.log` pattern)

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Ignored {
		t.Fatal("expected app.log to be gitignored, got false")
	}
}
```
