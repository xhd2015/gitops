## Expected
- `CompareBranches` returns `BranchRelationBIsAncestorOfA` (b is ancestor of a, so a is ahead)
- `CommitsAheadA` is 1 (main has one commit not in a)

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
	if resp.Result.Relation != git.BranchRelationBIsAncestorOfA {
		t.Fatalf("expected BranchRelationBIsAncestorOfA, got %v", resp.Result.Relation)
	}
	if resp.Result.CommitsAheadA != 1 {
		t.Fatalf("expected CommitsAheadA=1, got %d", resp.Result.CommitsAheadA)
	}
	if resp.Result.CommitsAheadB != 0 {
		t.Fatalf("expected CommitsAheadB=0, got %d", resp.Result.CommitsAheadB)
	}
}
```
