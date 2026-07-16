## Expected

- `DiffCached` returns a non-nil `*model.CachedDiff` with no error.
- `len(Files) >= 1`.
- At least one FilePatch has OldPath or NewPath equal to / containing `file.txt`.
- That FilePatch has non-empty Kind.
- Unified / Raw content is not required in this phase.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error for staged file: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Diff == nil {
		t.Fatal("expected non-nil Diff for staged file, got nil")
	}
	if len(resp.Diff.Files) < 1 {
		t.Fatalf("expected len(Files) >= 1, got %d", len(resp.Diff.Files))
	}
	found := false
	for _, fp := range resp.Diff.Files {
		pathHit := fp.NewPath == "file.txt" || fp.OldPath == "file.txt" ||
			strings.Contains(fp.NewPath, "file.txt") || strings.Contains(fp.OldPath, "file.txt")
		if !pathHit {
			continue
		}
		found = true
		if fp.Kind == "" {
			t.Fatalf("FilePatch for file.txt has empty Kind: %#v", fp)
		}
		break
	}
	if !found {
		t.Fatalf("expected a FilePatch path mentioning file.txt, got Files=%#v", resp.Diff.Files)
	}
}
```
