## Expected
- `CompareBranches` returns `BranchRelationDiverged`
- `DiffFileCount` is 0 (same content on both sides)
- `MergeBase` is non-empty
- `CommitsAheadA` is 1 and `CommitsAheadB` is 1

```go
import (
	"testing"

	"github.com/xhd2015/gitops/git"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Result.Relation != git.BranchRelationDiverged {
		t.Fatalf("expected BranchRelationDiverged, got %v", resp.Result.Relation)
	}
	if resp.Result.DiffFileCount != 0 {
		t.Fatalf("expected DiffFileCount=0, got %d", resp.Result.DiffFileCount)
	}
	if resp.Result.MergeBase == "" {
		t.Fatal("expected non-empty MergeBase")
	}
	if resp.Result.CommitsAheadA != 1 {
		t.Fatalf("expected CommitsAheadA=1, got %d", resp.Result.CommitsAheadA)
	}
	if resp.Result.CommitsAheadB != 1 {
		t.Fatalf("expected CommitsAheadB=1, got %d", resp.Result.CommitsAheadB)
	}
}
```
