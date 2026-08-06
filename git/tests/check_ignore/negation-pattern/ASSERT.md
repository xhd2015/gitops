## Expected
- `CheckIgnore` returns `false` (`important.o` is un-ignored by `!important.o`)

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
		t.Fatal("expected important.o to NOT be gitignored (negation), got true")
	}
}
```
