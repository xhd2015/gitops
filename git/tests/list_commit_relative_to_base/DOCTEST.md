# ListCommitRelativeToBase second-parent recovery

## Version

0.0.2

# DSN (Domain Specific Notion)

A Git branch has a head, a base, and a first-parent history. A merge commit on
that history has a second parent representing the line merged into the branch.
`ListCommitRelativeToBase` and `FindDiffPoints` are the two public branch
history views. They must retain the first-parent history and recover the
first-parent path of a directly merged second parent, but must not recursively
widen into merges encountered on that recovered path.

## Decision tree

```text
commit graph
├── linear first-parent path
├── merge with recovered second-parent path
├── recovered path containing an inner merge
└── multiple merge paths that reconverge
```

## Test case index

| Path | Graph property | Expected result |
|---|---|---|
| `linear/` | no merge | only the ordinary first-parent commits |
| `stale-local-merge/` | merge has stale P1 and remote P2 | all P2 first-parent commits are recovered once |
| `inner-merge-boundary/` | recovered path contains a merge | its own P2 is not recursively recovered |
| `reconverging-side-chains/` | two direct merges share history | every recovered commit occurs once |

## How to run

```sh
doctest vet ./git/tests/list_commit_relative_to_base
doctest test ./git/tests/list_commit_relative_to_base
```

```go
import (
    "testing"

    "github.com/xhd2015/doctest/session"
    gitopsgit "github.com/xhd2015/gitops/git"
)

type Request struct {
    Dir      string
    Head     string
    Base     string
    Expected []string
}

type Response struct {
    RelativeHashes []string
    DiffPointHashes []string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
    exists, merged, commits, err := gitopsgit.ListCommitRelativeToBase(req.Dir, req.Head, req.Base)
    if err != nil {
        return nil, err
    }
    if !exists {
        t.Fatal("expected graph refs to exist")
    }
    if merged {
        t.Fatal("expected head not to be merged into base")
    }
    hashes := make([]string, 0, len(commits))
    for _, commit := range commits {
        hashes = append(hashes, commit.Hash)
    }
    _, _, _, _, _, diffPointHashes, err := gitopsgit.FindDiffPoints(req.Dir, req.Head, req.Base, "")
    if err != nil {
        return nil, err
    }
    return &Response{RelativeHashes: hashes, DiffPointHashes: diffPointHashes}, nil
}
```
