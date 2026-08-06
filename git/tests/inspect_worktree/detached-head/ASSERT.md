## Expected

- `Branch` is `(detached)`
- `CommitShort` is 7 hex characters
- `CommitMessage` is `init`
- `IsClean` is true

```go
import (
	"regexp"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inspect := resp.Inspect
	if inspect.Branch != "(detached)" {
		t.Fatalf("Branch = %q, want (detached)", inspect.Branch)
	}
	if !regexp.MustCompile(`^[0-9a-f]{7}$`).MatchString(inspect.CommitShort) {
		t.Fatalf("CommitShort = %q, want 7-char hex hash", inspect.CommitShort)
	}
	if inspect.CommitMessage != "init" {
		t.Fatalf("CommitMessage = %q, want init", inspect.CommitMessage)
	}
	if !inspect.IsClean {
		t.Fatalf("expected clean detached worktree, got %+v", inspect)
	}
}
```