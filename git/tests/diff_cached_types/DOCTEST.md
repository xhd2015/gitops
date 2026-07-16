# DiffCached types — CachedDiff + DiffCachedParseError

## Version

0.0.2

P1 type surface for staged unified-diff parsing: `model.CachedDiff`,
`model.FilePatch`, `model.Hunk`, and `git.DiffCachedParseError`.

This tree does **not** call `DiffCached` or change its signature (P2).

# DSN (Domain Specific Notion)

**Participants**

- **CachedDiff** — in-memory model of a parsed staged patch: file list plus
  optional raw text.
- **FilePatch** — one file’s change: old/new path, kind, binary flag, hunks.
- **Hunk** — one hunk header plus body lines.
- **DiffCachedParseError** — typed error when parse of `git diff --cached`
  output fails; carries `Dir`, `Raw`, and wrapped `Err`.

**Behaviors**

- Zero-value and field-populated `CachedDiff` / `FilePatch` / `Hunk` are
  constructible and field-readable.
- `DiffCachedParseError.Error()` is non-empty and mentions parse or diff.
- `Unwrap()` returns the underlying `Err`.
- `errors.As` recovers `Dir` and `Raw`.

## Decision Tree

Split on **subject under test** (types vs parse-error contract):

```
diff_cached_types/
├── cached-diff-fields/   # model.CachedDiff / FilePatch / Hunk field surface
└── parse-error-raw/      # git.DiffCachedParseError Error/Unwrap/As
```

## Test Case Index

| # | Path | Preconditions | Expected |
|---|------|---------------|----------|
| 1 | `cached-diff-fields/` | construct CachedDiff with one FilePatch + Hunk + Raw | fields readable; zero-value Files nil/empty; no FileCount helper |
| 2 | `parse-error-raw/` | construct DiffCachedParseError with Dir, Raw, Err | Error() non-empty mentions parse/diff; Unwrap=Err; As recovers Raw+Dir |

## How to Run

```sh
cd external/gitops-master-2026-07-16
doctest vet ./git/tests/diff_cached_types
doctest test -v ./git/tests/diff_cached_types
```

```go
import (
	"errors"
	"fmt"
	"testing"

	"github.com/xhd2015/gitops/git"
	"github.com/xhd2015/gitops/model"
)

// Mode selects which type surface Run exercises.
// "cached-diff-fields" | "parse-error-raw"
type Request struct {
	Mode string

	// parse-error-raw inputs
	Dir      string
	Raw      string
	CauseMsg string

	// cached-diff-fields inputs (sample patch)
	SampleRaw      string
	SampleOldPath  string
	SampleNewPath  string
	SampleKind     string
	SampleBinary   bool
	SampleHunkHdr  string
	SampleHunkLine string
}

type Response struct {
	// cached-diff-fields
	ZeroFilesLen int
	FilesLen     int
	Raw          string
	OldPath      string
	NewPath      string
	Kind         string
	Binary       bool
	HunksLen     int
	HunkHeader   string
	HunkLines    []string

	// parse-error-raw
	ErrMsg       string
	UnwrappedMsg string
	AsOK         bool
	AsDir        string
	AsRaw        string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	switch req.Mode {
	case "cached-diff-fields":
		var zero model.CachedDiff
		populated := model.CachedDiff{
			Raw: req.SampleRaw,
			Files: []model.FilePatch{
				{
					OldPath: req.SampleOldPath,
					NewPath: req.SampleNewPath,
					Kind:    req.SampleKind,
					Binary:  req.SampleBinary,
					Hunks: []model.Hunk{
						{
							Header: req.SampleHunkHdr,
							Lines:  []string{req.SampleHunkLine},
						},
					},
				},
			},
		}
		resp := &Response{
			ZeroFilesLen: len(zero.Files),
			FilesLen:     len(populated.Files),
			Raw:          populated.Raw,
		}
		if len(populated.Files) > 0 {
			fp := populated.Files[0]
			resp.OldPath = fp.OldPath
			resp.NewPath = fp.NewPath
			resp.Kind = fp.Kind
			resp.Binary = fp.Binary
			resp.HunksLen = len(fp.Hunks)
			if len(fp.Hunks) > 0 {
				resp.HunkHeader = fp.Hunks[0].Header
				resp.HunkLines = append([]string(nil), fp.Hunks[0].Lines...)
			}
		}
		return resp, nil

	case "parse-error-raw":
		cause := errors.New(req.CauseMsg)
		perr := &git.DiffCachedParseError{
			Dir: req.Dir,
			Raw: req.Raw,
			Err: cause,
		}
		var as *git.DiffCachedParseError
		ok := errors.As(perr, &as)
		resp := &Response{
			ErrMsg: perr.Error(),
			AsOK:   ok,
		}
		if u := errors.Unwrap(perr); u != nil {
			resp.UnwrappedMsg = u.Error()
		}
		if ok && as != nil {
			resp.AsDir = as.Dir
			resp.AsRaw = as.Raw
		}
		return resp, perr

	default:
		return nil, fmt.Errorf("unknown mode %q", req.Mode)
	}
}
```
