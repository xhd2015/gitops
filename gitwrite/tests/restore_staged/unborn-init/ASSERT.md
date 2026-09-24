## Expected

- `RestoreStaged` succeeds (must not `could not resolve 'HEAD'`).
- `drop.bin` is no longer staged.
- `keep.txt` is still staged.
- `drop.bin` still exists on disk.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("RestoreStaged on unborn HEAD: %v", err)
	}
	hasDrop, hasKeep := false, false
	for _, f := range resp.StagedAfter {
		if f == "drop.bin" {
			hasDrop = true
		}
		if f == "keep.txt" {
			hasKeep = true
		}
	}
	if hasDrop {
		t.Fatalf("expected drop.bin unstaged, staged=%v", resp.StagedAfter)
	}
	if !hasKeep {
		t.Fatalf("expected keep.txt still staged, staged=%v", resp.StagedAfter)
	}
	if !fileExists(filepath.Join(req.Dir, "drop.bin")) {
		t.Fatal("expected drop.bin to still exist on disk")
	}
}
```
