## Expected
- `CompareBranches` returns an error indicating the directory is not found

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for nonexistent dir, got nil")
	}
	if !strings.Contains(err.Error(), "no such file") && !strings.Contains(err.Error(), "fatal:") && !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("expected error to indicate directory not found, got: %s", err.Error())
	}
}
```
