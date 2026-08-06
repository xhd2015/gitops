## Expected

- `IsSubmodule` returns no error.
- `IsSubmodule(dir, "ext")` returns `false` — `ext/` is a nested separate repo, not
  a tracked submodule of `dir`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.Submodule {
		t.Fatal("expected ext/ to NOT be a submodule (nested separate repo), got true")
	}
}
```
