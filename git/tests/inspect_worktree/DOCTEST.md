# InspectWorktree Doctests

Doc-style tests for `git.InspectWorktree(dir)` — live branch, HEAD commit,
and porcelain change counts for a directory.

# DSN (Domain Specific Notion)

**Participants**

- **git repository** — temp dir created per leaf; may be absent (not-a-repo).
- **InspectWorktree** — reads branch (`git branch --show-current`), HEAD commit
  via `GetCommit`, and `git status --porcelain=v1` for change-type counts.
- **Porcelain classifier** — `classifyPorcelainLine` maps XY status pairs to
  added / changed / renamed / deleted buckets.

**Behaviors**

- Non-git dir → `IsRepo=false`, empty branch/commit, `IsClean=true`, zero counts.
- Clean repo → branch name, 7-char commit hash, subject message, `IsClean=true`.
- Dirty repo → per-type counts; untracked (`??`) counts as **added**.
- Detached HEAD → `Branch="(detached)"`, commit still populated.
- Porcelain rules: `??`→added; `R`→renamed; `D`→deleted; `A`→added; `M`→changed.
- Untracked directories expand **added** count via `git ls-files --others`.

## Version

0.0.2

## Decision Tree

```
[InspectWorktree(dir)]
 |
 +-- not-a-repo/              (LEAF)  plain directory, no .git
 +-- clean-repo/              (LEAF)  single commit, clean worktree
 +-- dirty-all-types/        (LEAF)  added + changed + renamed + deleted
 +-- detached-head/          (LEAF)  git checkout --detach HEAD
 +-- untracked-dir/          (LEAF)  untracked dir expands added file count
 +-- porcelain-classifier/
      +-- representative/     (LEAF)  classifyPorcelainLine table
```

## Test Index

| # | Leaf | Description |
|---|------|-------------|
| 1 | `not-a-repo` | Non-git dir → `IsRepo=false`, clean zero state |
| 2 | `clean-repo` | `master`, 7-char hash, message `init`, clean |
| 3 | `dirty-all-types` | 1 added, 1 changed, 1 renamed, 1 deleted |
| 4 | `detached-head` | Branch `(detached)`, commit present, clean |
| 5 | `untracked-dir` | Untracked dir counts files inside (`Added=4`) |
| 6 | `porcelain-classifier/representative` | XY → kind mapping table |

## How to Run

```sh
cd gitops-adhoc
doctest vet ./git/tests/inspect_worktree
doctest test ./git/tests/inspect_worktree/...
```

```go
import (
	"testing"

	"github.com/xhd2015/gitops/git"
)

type PorcelainCase struct {
	XY   string
	Want string
}

type Request struct {
	Dir            string
	PorcelainCases []PorcelainCase
}

type Response struct {
	Inspect       *git.WorktreeInspect
	PorcelainKind string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if len(req.PorcelainCases) > 0 {
		return &Response{}, nil
	}
	inspect, err := git.InspectWorktree(req.Dir)
	if err != nil {
		return nil, err
	}
	return &Response{Inspect: inspect}, nil
}
```