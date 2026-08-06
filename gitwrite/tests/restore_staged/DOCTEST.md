# RestoreStaged Test Case Tree

## Version

0.0.2

Run with:
```sh
doctest test ./ -v
```

## DSN (Domain Specific Notion)

The **caller** hands `RestoreStaged(dir, paths...)` a directory inside a git work tree
and one or more file paths. The **helper** runs `git restore --staged -- <paths>` in
`dir`, removing those paths from the git index (staging area) while leaving the
working copy untouched.

## Decision Tree

The most significant factor is **how many files are unstaged**; the partial case
verifies only the specified files are affected.

```
RestoreStaged(dir, paths...)
├── single-file/        # unstage one file, verify gone from index, file still on disk
├── multiple-files/     # unstage two files, verify both gone
└── partial-subset/     # unstage one of two staged files, verify only that one is gone
```

## Test Case Index

| # | Path | Preconditions | Expected |
|---|---|---|---|
| 1 | `single-file/` | `a.txt` staged | `a.txt` not staged, `a.txt` still on disk |
| 2 | `multiple-files/` | `a.txt` and `b.txt` staged | both not staged, both still on disk |
| 3 | `partial-subset/` | `a.txt` and `b.txt` staged, restore only `a.txt` | `a.txt` not staged, `b.txt` still staged, both on disk |

## How to Run

```sh
doctest vet ./external/gitops-master-2026-07-02/gitwrite/tests/restore_staged
doctest test ./external/gitops-master-2026-07-02/gitwrite/tests/restore_staged
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

type Request struct {
	Dir   string
	Paths []string
}

type Response struct {
	// StagedAfter is the set of staged files after RestoreStaged completes.
	StagedAfter []string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	err := RestoreStaged(req.Dir, req.Paths...)
	if err != nil {
		return nil, err
	}
	files, err := getStagedFileNames(req.Dir)
	if err != nil {
		return nil, err
	}
	return &Response{StagedAfter: files}, nil
}
```
