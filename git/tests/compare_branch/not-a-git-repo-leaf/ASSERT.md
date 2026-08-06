## Expected
- `CompareBranches` returns an error indicating this is not a git repository

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for non-git directory, got nil")
	}
	if !strings.Contains(err.Error(), "not a git") && !strings.Contains(err.Error(), "fatal:") {
		t.Fatalf("expected error to indicate not a git repository, got: %s", err.Error())
	}
}
```
