## Expected
- `CheckIgnore` returns `false` (non-existent paths are not ignored)
- `git check-ignore` returns exit code 1 when path does not exist

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
		t.Fatal("expected non-existent path to NOT be gitignored, got true")
	}
}
```
