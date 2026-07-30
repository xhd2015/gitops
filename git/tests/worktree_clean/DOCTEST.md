# Dual Clean Doctests (diff-clean vs porcelain-clean)

Doc-style tests for distinctly named clean APIs on
`github.com/xhd2015/gitops/git/worktree`:

- **IsDiffClean** — `git diff --quiet HEAD` semantics (untracked may still be clean)
- **IsPorcelainClean** — empty `git status --porcelain` (untracked ⇒ dirty)

# DSN (Domain Specific Notion)

**Participants**

- **Worktree directory** — temp git checkout under test.
- **IsDiffClean** — true when tracked-tree diff vs HEAD is empty (diff exit 0).
- **IsPorcelainClean** — true when porcelain status output is empty.

**Behaviors**

- Clean repo (no untracked, no modifications): both true.
- Untracked-only file: porcelain-clean false; diff-clean true.
- Modified tracked file: both false.
- APIs return `(bool, error)` and are named distinctly (not a single overloaded clean).

## Version

0.0.2

## Decision Tree

```
worktree_clean
├── clean-repo/           (LEAF) both true
├── untracked-only/       (LEAF) diff true, porcelain false
└── modified-tracked/     (LEAF) both false
```

## Test Index

| # | Leaf | Description |
|---|------|-------------|
| 1 | `clean-repo` | No changes → DiffClean=true, PorcelainClean=true |
| 2 | `untracked-only` | Untracked file → DiffClean=true, PorcelainClean=false |
| 3 | `modified-tracked` | Modified tracked file → both false |

## How to Run

```sh
cd external/gitops-master-2026-07-30
doctest vet ./git/tests/worktree_clean
doctest test ./git/tests/worktree_clean/...
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/gitops/git/worktree"
)

type Request struct {
	Dir string
}

type Response struct {
	DiffClean      bool
	PorcelainClean bool
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	diffOK, err := worktree.IsDiffClean(req.Dir)
	if err != nil {
		return nil, err
	}
	porcOK, err := worktree.IsPorcelainClean(req.Dir)
	if err != nil {
		return nil, err
	}
	return &Response{DiffClean: diffOK, PorcelainClean: porcOK}, nil
}
```
