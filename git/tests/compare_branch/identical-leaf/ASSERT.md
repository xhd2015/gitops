## Expected
- `CompareBranches` returns `BranchRelationSame`

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/gitops/git"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Result.Relation != git.BranchRelationSame {
		t.Fatalf("expected BranchRelationSame, got %v", resp.Result.Relation)
	}
}
```
