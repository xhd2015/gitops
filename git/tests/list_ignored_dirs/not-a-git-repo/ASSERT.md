## Expected

- `ListIgnoredDirs` returns no error.
- The result is empty (git-based ignore listing is inapplicable when dir is not a
  git repo).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Dirs) != 0 {
		t.Fatalf("ListIgnoredDirs = %v, want empty (not a git repo)", resp.Dirs)
	}
}
```
