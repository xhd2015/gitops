# GetStagedFiles Test Case Tree

## Version

0.0.2

Run with:
```sh
doctest test ./ -v
```

## DSN (Domain Specific Notion)

The **caller** hands `GetStagedFiles(dir)` a directory inside a git work tree. The
**helper** runs `git diff --cached --name-only --diff-filter=ACMRT --` in `dir`
and returns the list of staged file paths relative to `dir`. Files with status `D`
(deleted) are excluded from the result because the diff filter excludes them.

## Decision Tree

The most significant factor is **how many staged files exist**; edge cases cover
deleted exclusions and subdirectory paths.

```
GetStagedFiles(dir)
├── empty-index/          # no staged files → []
├── single-file/          # one staged file → [file.txt]
├── multiple-files/       # multiple staged files → all listed
├── deleted-excluded/     # deleted staged file excluded by diff-filter
└── subdirectory-file/    # file in subdir returned with correct relative path
```

## Test Case Index

| # | Path | Preconditions | Expected |
|---|---|---|---|
| 1 | `empty-index/` | Fresh repo with initial commit, nothing staged | `[]` |
| 2 | `single-file/` | `file.txt` staged | `["file.txt"]` |
| 3 | `multiple-files/` | `a.txt` and `b.txt` staged | `["a.txt", "b.txt"]` |
| 4 | `deleted-excluded/` | `remove-me.txt` committed then deleted and staged | `[]` (excluded) |
| 5 | `subdirectory-file/` | `sub/pkg.go` staged | `["sub/pkg.go"]` |

## How to Run

```sh
doctest vet ./external/gitops-master-2026-07-02/git/tests/get_staged_files
doctest test ./external/gitops-master-2026-07-02/git/tests/get_staged_files
```

```go
import (
	"testing"
)

type Request struct {
	Dir string
}

type Response struct {
	Files []string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	files, err := GetStagedFiles(req.Dir)
	if err != nil {
		return nil, err
	}
	return &Response{Files: files}, nil
}
```
