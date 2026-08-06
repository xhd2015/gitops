## Expected

- `IsSubmodule` returns no error.
- `IsSubmodule(dir, "ext")` returns `true` — `ext/` is a tracked submodule gitlink
  (mode 160000) of `dir`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Submodule {
		t.Fatal("expected ext/ to be a tracked submodule, got false")
	}
}
```
