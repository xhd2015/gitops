## Expected
- `CheckIgnore` returns `false` (`main.go` does not match `*.log`)

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.Ignored {
		t.Fatal("expected main.go to NOT be gitignored, got true")
	}
}
```
