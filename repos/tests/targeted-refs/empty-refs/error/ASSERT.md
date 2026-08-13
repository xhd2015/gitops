## Expected

- `Run` returns a non-nil error.
- Error is not an empty string.

## Side Effects

- No dest-ref inventory is required (fetch must not succeed).

## Errors

- `FetchTargetedRefs` with no refs returns an error.

## Exit Code

- Harness exit is zero (`Assert` accepts the product error).

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
	_ = resp
	if err == nil {
		t.Fatal("FetchTargetedRefs with no refs: expected error")
	}
	if strings.TrimSpace(err.Error()) == "" {
		t.Fatal("error must have a message")
	}
}
```
