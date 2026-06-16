## Expected
- `CompareBranches` returns `BranchRelationAIsAncestorOfB` (a is ancestor of b, so b is ahead)
- `CommitsAheadB` is 1 (main has one commit not in a)

```go
import (
	"testing"

	"github.com/xhd2015/gitops/git"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Result.Relation != git.BranchRelationAIsAncestorOfB {
		t.Fatalf("expected BranchRelationAIsAncestorOfB, got %v", resp.Result.Relation)
	}
	if resp.Result.CommitsAheadB != 1 {
		t.Fatalf("expected CommitsAheadB=1, got %d", resp.Result.CommitsAheadB)
	}
	if resp.Result.CommitsAheadA != 0 {
		t.Fatalf("expected CommitsAheadA=0, got %d", resp.Result.CommitsAheadA)
	}
}
```
