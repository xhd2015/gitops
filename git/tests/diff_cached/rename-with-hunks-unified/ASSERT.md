## Expected

- `ParseCachedDiff` returns a non-nil `*model.CachedDiff` with no error.
- One `FilePatch` with `Kind == "rename"`, paths `old.go` → `new.go`, and
  `len(Hunks) >= 1`.
- `Unified()` includes rename meta **and** content section:
  - `rename from old.go` / `rename to new.go`
  - `--- a/old.go` and `+++ b/new.go`
  - at least one `@@` hunk header
  - recognizable hunk body (`+new line` or `-old line`)

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error for rename+hunks parse: %v", err)
	}
	if resp == nil || resp.Diff == nil {
		t.Fatal("expected non-nil Diff for rename+hunks")
	}

	found := false
	for _, fp := range resp.Diff.Files {
		if fp.Kind != "rename" {
			continue
		}
		if fp.OldPath != "old.go" || fp.NewPath != "new.go" {
			continue
		}
		found = true
		if len(fp.Hunks) == 0 {
			t.Fatal("rename+content expected non-empty Hunks")
		}
	}
	if !found {
		t.Fatalf("expected FilePatch Kind=rename old.go→new.go; got %#v", resp.Diff.Files)
	}

	u := resp.Diff.Unified()
	if u == "" {
		t.Fatal("Unified() returned empty string")
	}
	for _, want := range []string{
		"diff --git a/old.go b/new.go",
		"rename from old.go",
		"rename to new.go",
		"--- a/old.go",
		"+++ b/new.go",
		"@@",
	} {
		if !strings.Contains(u, want) {
			t.Fatalf("rename+hunks Unified() missing %q; got:\n%s", want, u)
		}
	}
	if !strings.Contains(u, "+new line") && !strings.Contains(u, "-old line") {
		t.Fatalf("rename+hunks Unified() missing hunk body markers; got:\n%s", u)
	}
}
```
