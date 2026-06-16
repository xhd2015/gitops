## Expected
- `CompareBranches` returns an error
- The error message contains the invalid ref name

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for invalid ref, got nil")
	}
	if !strings.Contains(err.Error(), "nonexistent-branch") {
		t.Fatalf("expected error to mention 'nonexistent-branch', got: %s", err.Error())
	}
}
```
