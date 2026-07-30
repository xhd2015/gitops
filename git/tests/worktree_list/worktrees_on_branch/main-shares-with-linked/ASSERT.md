## Expected

- Length 2
- Includes main and linked
- Both Branch `feature`
- One entry IsMain true (main)

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Entries) != 2 {
		t.Fatalf("WorktreesOnBranch(feature) len = %d, want 2; paths=%v", len(resp.Entries), entryPaths(resp.Entries))
	}
	if !containsPath(t, resp.Entries, req.MainPath) {
		t.Fatalf("missing main %q; paths=%v", req.MainPath, entryPaths(resp.Entries))
	}
	if !containsPath(t, resp.Entries, req.LinkedPath) {
		t.Fatalf("missing linked %q; paths=%v", req.LinkedPath, entryPaths(resp.Entries))
	}
	mainE, ok := findByPath(t, resp.Entries, req.MainPath)
	if !ok || !mainE.IsMain {
		t.Fatalf("main IsMain false or missing: %+v", mainE)
	}
	if mainE.Branch != "feature" {
		t.Fatalf("main Branch = %q, want feature", mainE.Branch)
	}
	linkedE, ok := findByPath(t, resp.Entries, req.LinkedPath)
	if !ok {
		t.Fatal("linked missing")
	}
	if linkedE.IsMain {
		t.Fatal("linked IsMain=true, want false")
	}
	if linkedE.Branch != "feature" {
		t.Fatalf("linked Branch = %q, want feature", linkedE.Branch)
	}
}
```
