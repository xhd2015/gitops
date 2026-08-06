# Scenario

**Feature**: model.CachedDiff, FilePatch, and Hunk field surface

```
# zero-value CachedDiff has empty Files
var zero model.CachedDiff -> len(Files)==0

# populated CachedDiff carries Raw + one FilePatch with Hunk
CachedDiff{Raw, Files:[{OldPath,NewPath,Kind,Binary,Hunks:[{Header,Lines}]}]}
  -> fields readable via Run response
```

## Preconditions

- Types `model.CachedDiff`, `model.FilePatch`, `model.Hunk` exist with the
  locked field names.
- No FileCount (or similar) helper is required yet.

## Steps

1. Set `req.Mode` to `cached-diff-fields`.
2. Provide sample path/kind/hunk/raw values on the request.
3. Run constructs a zero-value and a populated `CachedDiff`.

## Context

- Kind values are free strings for P1 (`modify` used as sample).
- Binary false for the text sample.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "cached-diff-fields"
	req.SampleRaw = "diff --git a/file.txt b/file.txt\n"
	req.SampleOldPath = "file.txt"
	req.SampleNewPath = "file.txt"
	req.SampleKind = "modify"
	req.SampleBinary = false
	req.SampleHunkHdr = "@@ -1 +1 @@"
	req.SampleHunkLine = "+hello"
	return nil
}
```
