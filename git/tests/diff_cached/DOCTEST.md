# DiffCached structured I/O + CachedDiff render helpers

## Version

0.0.2

Run with:
```sh
doctest test ./ -v
```

# DSN (Domain Specific Notion)

**Participants**

- **caller** — hands a work-tree directory (or raw patch text for the parse helper),
  or holds a `*model.CachedDiff` and asks for render helpers.
- **DiffCached** — runs `git -C dir diff --cached`, then parses stdout into
  `*model.CachedDiff` (Files with Kind/paths) or returns a typed parse error.
- **ParseCachedDiff** — pure parser over a raw unified-diff string; same success
  and parse-failure contract as DiffCached without invoking git.
- **CachedDiff** — in-memory staged patch model; exposes `FileCount`, `Unified`,
  and `UnifiedTruncated` for consumers (e.g. commit_msg truncation).
- **DiffCachedParseError** — typed failure when non-empty stdout cannot be
  parsed; carries `Raw` (full input/stdout) and optional `Dir`.
- **git repository / plain directory** — live preconditions for staging-area
  and non-repo error leaves.

**Behaviors**

- Empty staged patch (git ok, empty stdout) → `(nil, nil)`.
- Staged content → non-nil `*CachedDiff` with `Files` populated (paths and Kind set).
- Not a git repo / fatal git failure → ordinary `(nil, err)` (not parse-error).
- Non-empty unparseable text → `(nil, *DiffCachedParseError)` with `Raw` preserved.
- `FileCount()` → number of file patches; nil receiver → 0.
- `Unified()` → full unified-diff text for the staged patch; nil receiver → `""`.
- **Pure rename** (`Kind=rename`, non-binary, empty `Hunks`): `Unified` emits
  only `diff --git a/<old> b/<new>`, `rename from <old>`, `rename to <new>` —
  no `---/+++` file headers and no `@@` hunks.
- **Rename + content** (`Kind=rename` with non-empty `Hunks`): rename meta plus
  `---/+++` and hunk bodies (same shape as git for partial renames).
- `UnifiedTruncated(maxLinesPerFile)` → per file section, keep first
  `maxLinesPerFile` lines; if more, append `\n...(%d more lines omitted)` with
  remaining line count; nil receiver → `""`.

## Decision Tree

Most significant factor is **outcome path** (empty vs staged content vs git
failure vs parse failure vs CachedDiff method surface). Live leaves exercise
`DiffCached(dir)`; the parse leaf exercises exported `ParseCachedDiff(raw)`;
method leaves call helpers on the returned (or nil) `*model.CachedDiff`.

```
DiffCached / ParseCachedDiff / CachedDiff methods
├── empty-index/                 # nothing staged → (nil, nil)
├── after-unstage-all/           # stage then restore --staged → (nil, nil)
├── single-staged-file/          # one staged file → non-nil Files, path + Kind
├── not-a-git-repo/              # plain temp dir → ordinary error
├── parse-malformed/             # ParseCachedDiff(bad raw) → *DiffCachedParseError
├── file-count/                  # two staged files → FileCount()==2
├── unified-roundtrip-ish/       # staged small file → Unified has diff --git + path
├── unified-truncated-long-file/ # long staged file → UnifiedTruncated(24) omits tail
├── pure-rename-unified/         # live 100% rename → Unified: rename meta only (no ---/+++, no @@)
├── rename-with-hunks-unified/   # parse rename+edit → Unified keeps rename meta + ---/+++ + @@
└── nil-receiver-methods/        # nil *CachedDiff → FileCount 0, Unified/Truncated ""
```

## Test Case Index

| # | Path | Preconditions | Expected |
|---|---|---|---|
| 1 | `empty-index/` | Fresh repo with initial commit, nothing staged | `Diff == nil`, err nil |
| 2 | `after-unstage-all/` | Stage then `git restore --staged -- .` | `Diff == nil`, err nil |
| 3 | `single-staged-file/` | `file.txt` staged with content | non-nil Diff; `len(Files)>=1`; path includes `file.txt`; Kind set |
| 4 | `not-a-git-repo/` | Temp dir with no `.git` | non-nil ordinary error; not required to be ParseError |
| 5 | `parse-malformed/` | Malformed raw via `ParseCachedDiff` | `(nil, *DiffCachedParseError)` with `Raw` = input |
| 6 | `file-count/` | Two distinct staged files | `FileCount() == 2` |
| 7 | `unified-roundtrip-ish/` | One small staged file (`note.txt`) | `Unified()` contains `diff --git` and `note.txt` |
| 8 | `unified-truncated-long-file/` | One staged file whose unified section is >24 lines | `UnifiedTruncated(24)` has omit marker; shorter than `Unified()`; late body absent |
| 9 | `pure-rename-unified/` | Live `git mv` 100% rename (`oldname.txt`→`newname.txt`) | `Unified()` has rename from/to; no `@@`; no `---/+++` |
| 10 | `rename-with-hunks-unified/` | Mode `parse` rename+content raw (`old.go`→`new.go` + hunk) | `Unified()` has rename meta + `---/+++` + `@@` + body |
| 11 | `nil-receiver-methods/` | Mode `nil-methods` (no git call) | nil receiver: FileCount 0, Unified `""`, UnifiedTruncated `""` |

## How to Run

```sh
cd external/gitops-master-2026-07-22
doctest vet ./git/tests/diff_cached
doctest test -v ./git/tests/diff_cached
```

```go
import (
	"testing"

	"github.com/xhd2015/gitops/git"
	"github.com/xhd2015/gitops/model"
)

// Mode selects the entrypoint under test.
// "" or "diff-cached" → git.DiffCached(req.Dir)
// "parse" → git.ParseCachedDiff(req.Raw)
// "nil-methods" → return Diff=nil without calling git (method-nil contract)
type Request struct {
	Dir  string
	Mode string
	Raw  string
}

type Response struct {
	Diff *model.CachedDiff
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Mode == "nil-methods" {
		return &Response{Diff: nil}, nil
	}
	if req.Mode == "parse" {
		diff, err := git.ParseCachedDiff(req.Raw)
		if err != nil {
			return &Response{Diff: diff}, err
		}
		return &Response{Diff: diff}, nil
	}
	diff, err := git.DiffCached(req.Dir)
	if err != nil {
		return &Response{Diff: diff}, err
	}
	return &Response{Diff: diff}, nil
}
```
