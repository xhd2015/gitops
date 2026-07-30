## Expected

- `WorktreesOnBranch(master)` is exactly main (len=1)
- Detached path never appears for `master` or `feature`
- No error

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	onMaster := resp.ByBranch["master"]
	if len(onMaster) != 1 {
		t.Fatalf("WorktreesOnBranch(master) len = %d, want 1; paths=%v", len(onMaster), entryPaths(onMaster))
	}
	if !samePath(t, onMaster[0].Path, req.MainPath) {
		t.Fatalf("master path = %q, want main %q", onMaster[0].Path, req.MainPath)
	}
	if containsPath(t, onMaster, req.LinkedPath) {
		t.Fatal("detached path must not match named branch master")
	}
	onFeature := resp.ByBranch["feature"]
	if containsPath(t, onFeature, req.LinkedPath) {
		t.Fatal("detached path must not match named branch feature")
	}
	if len(onFeature) != 0 {
		t.Fatalf("WorktreesOnBranch(feature) len = %d, want 0", len(onFeature))
	}
}
```
