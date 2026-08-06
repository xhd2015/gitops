## Expected

- `DiffCached` returns a non-nil `*model.CachedDiff` with no error.
- At least one `FilePatch` has `Kind == "rename"`, empty `Hunks`, and paths
  covering `oldname.txt` → `newname.txt`.
- `Unified()` includes:
  - `diff --git a/oldname.txt b/newname.txt` (or equivalent a/ b/ form)
  - `rename from oldname.txt`
  - `rename to newname.txt`
- `Unified()` does **not** include:
  - any `@@` hunk header
  - any `--- ` / `+++ ` file headers (path-only rename has no content section)

## Expected Output

```text
diff --git a/oldname.txt b/newname.txt
rename from oldname.txt
rename to newname.txt
```

(Similarity-index lines may be present or absent in re-render; not required.)

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error for pure rename: %v", err)
	}
	if resp == nil || resp.Diff == nil {
		t.Fatal("expected non-nil Diff for pure rename")
	}

	foundRename := false
	for _, fp := range resp.Diff.Files {
		if fp.Kind != "rename" {
			continue
		}
		if fp.OldPath != "oldname.txt" || fp.NewPath != "newname.txt" {
			continue
		}
		foundRename = true
		if fp.Binary {
			t.Fatal("pure rename must not be Binary")
		}
		if len(fp.Hunks) != 0 {
			t.Fatalf("pure rename expected empty Hunks, got %d: %#v", len(fp.Hunks), fp.Hunks)
		}
	}
	if !foundRename {
		t.Fatalf("expected FilePatch Kind=rename oldname.txt→newname.txt; got %#v", resp.Diff.Files)
	}

	u := resp.Diff.Unified()
	if u == "" {
		t.Fatal("Unified() returned empty string")
	}
	if !strings.Contains(u, "diff --git a/oldname.txt b/newname.txt") {
		t.Fatalf("Unified() missing pure-rename diff --git header; got:\n%s", u)
	}
	if !strings.Contains(u, "rename from oldname.txt") {
		t.Fatalf("Unified() missing %q; got:\n%s", "rename from oldname.txt", u)
	}
	if !strings.Contains(u, "rename to newname.txt") {
		t.Fatalf("Unified() missing %q; got:\n%s", "rename to newname.txt", u)
	}
	if strings.Contains(u, "@@") {
		t.Fatalf("pure-rename Unified() must not include @@ hunks; got:\n%s", u)
	}
	// Path-only rename: no content file headers (current render wrongly adds these).
	for _, line := range strings.Split(u, "\n") {
		if strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ ") {
			t.Fatalf("pure-rename Unified() must not include ---/+++ headers; got line %q in:\n%s", line, u)
		}
	}
}
```
