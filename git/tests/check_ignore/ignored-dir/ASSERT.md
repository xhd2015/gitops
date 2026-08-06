## Expected
- `CheckIgnore` returns `true` (`build/` matches `build/` pattern in .gitignore)

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
		t.Fatal("expected build/ to be gitignored, got false")
	}
}
```
