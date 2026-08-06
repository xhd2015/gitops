## Expected

- Zero-value `CachedDiff` has `Files` length 0 (no helper API required).
- Populated `CachedDiff` exposes:
  - `Raw` matching the sample unified-diff prefix
  - one `FilePatch` with OldPath/NewPath/Kind/Binary
  - one `Hunk` with Header and Lines

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.ZeroFilesLen != 0 {
		t.Fatalf("zero-value CachedDiff Files len = %d, want 0", resp.ZeroFilesLen)
	}
	if resp.FilesLen != 1 {
		t.Fatalf("populated Files len = %d, want 1", resp.FilesLen)
	}
	if resp.Raw != req.SampleRaw {
		t.Fatalf("Raw = %q, want %q", resp.Raw, req.SampleRaw)
	}
	if resp.OldPath != req.SampleOldPath {
		t.Fatalf("OldPath = %q, want %q", resp.OldPath, req.SampleOldPath)
	}
	if resp.NewPath != req.SampleNewPath {
		t.Fatalf("NewPath = %q, want %q", resp.NewPath, req.SampleNewPath)
	}
	if resp.Kind != req.SampleKind {
		t.Fatalf("Kind = %q, want %q", resp.Kind, req.SampleKind)
	}
	if resp.Binary != req.SampleBinary {
		t.Fatalf("Binary = %v, want %v", resp.Binary, req.SampleBinary)
	}
	if resp.HunksLen != 1 {
		t.Fatalf("Hunks len = %d, want 1", resp.HunksLen)
	}
	if resp.HunkHeader != req.SampleHunkHdr {
		t.Fatalf("Hunk.Header = %q, want %q", resp.HunkHeader, req.SampleHunkHdr)
	}
	if len(resp.HunkLines) != 1 || resp.HunkLines[0] != req.SampleHunkLine {
		t.Fatalf("Hunk.Lines = %#v, want [%q]", resp.HunkLines, req.SampleHunkLine)
	}
}
```
