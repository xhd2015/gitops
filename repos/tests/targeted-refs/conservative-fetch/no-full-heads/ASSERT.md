## Expected

- `FetchTargetedRefs` succeeds and performs a fetch (`Fetched` true).
- Progress contains `--filter=blob:none`.
- Progress does **not** contain `+refs/heads/*`.

## Side Effects

- Fetch command is printed to the injected Progress writer.

## Errors

- None.

## Exit Code

- Exit code is zero.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
	if err != nil {
		t.Fatalf("FetchTargetedRefs: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if !resp.Fetched {
		t.Fatal("expected a real fetch so Progress records argv")
	}
	got := resp.Progress
	if !strings.Contains(got, "--filter=blob:none") {
		t.Fatalf("Progress missing --filter=blob:none:\n%s", got)
	}
	if strings.Contains(got, "+refs/heads/*") {
		t.Fatalf("Progress must not contain +refs/heads/* (full-heads fetch):\n%s", got)
	}
}
```
