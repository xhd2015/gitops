## Expected

- `FetchTargetedRefs` succeeds.
- Dest `refs/gitops/targets/{sha256("feature")}` exists.
- Dest `refs/gitops/targets/{sha256("master")}` does **not** exist.
- No `refs/remotes/origin/*` tracking refs are created.

## Side Effects

- Only the requested dest target is created under the default prefix.

## Errors

- None.

## Exit Code

- Exit code is zero.

```go
import (
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
	featureDest := expectedDestRef("", "feature")
	masterDest := expectedDestRef("", "master")
	if !hasRef(resp.AllRefs, featureDest) {
		t.Fatalf("missing requested dest %s; refs=%v", featureDest, resp.AllRefs)
	}
	if hasRef(resp.AllRefs, masterDest) {
		t.Fatalf("unused master dest %s was created; refs=%v", masterDest, resp.AllRefs)
	}
	if hasRefPrefix(resp.AllRefs, "refs/remotes/origin/") {
		t.Fatalf("narrow fetch must not create refs/remotes/origin/*; refs=%v", resp.AllRefs)
	}
	if _, ok := resp.Commits["feature"]; !ok || resp.Commits["feature"] == "" {
		t.Fatalf("Commits[feature] missing; got %#v", resp.Commits)
	}
}
```
