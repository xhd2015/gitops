## Expected

- `GetStagedFiles` returns both `a.txt` and `b.txt`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 2 {
		t.Fatalf("expected 2 staged files, got %d: %v", len(resp.Files), resp.Files)
	}
	hasA, hasB := false, false
	for _, f := range resp.Files {
		switch f {
		case "a.txt":
			hasA = true
		case "b.txt":
			hasB = true
		}
	}
	if !hasA {
		t.Fatal("expected a.txt in staged files")
	}
	if !hasB {
		t.Fatal("expected b.txt in staged files")
	}
}
```
