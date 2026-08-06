# IsSubmodule Test Case Tree

Run with:
```sh
doctest test ./ -v
```

## Version

0.0.2

## DSN (Domain Specific Notion)

The **caller** hands `IsSubmodule(dir, path)` a repository `dir` and a `path` relative
to it. The **helper** runs `git -C <dir> ls-files --stage -z -- <path>` and reports
`true` only when the path's index entry has mode `160000` — i.e. it is a tracked
submodule gitlink of `dir`. A nested directory that merely carries its own `.git`
but is NOT tracked by `dir` (a nested separate repo) returns `false`. When `dir` is
not a git repo, git is unavailable, or the path is not tracked at all, `IsSubmodule`
returns `false` and a nil error.

## Decision Tree

The most significant factor is **whether path is tracked by dir as a gitlink**; that
single factor separates a real submodule from a nested separate repo.

```
IsSubmodule(dir, path)
├── submodule-tracked/        # dir adds ext/ as a real submodule -> true
└── nested-separate-repo/     # ext/ has own .git, untracked by dir -> false
```

## How to Run

```sh
doctest vet ./gitops-dot/git/tests/is_submodule
doctest test ./gitops-dot/git/tests/is_submodule
```

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
)

type Request struct {
	Dir  string
	Path string
}

type Response struct {
	Submodule bool
}

// Run calls IsSubmodule(req.Dir, req.Path). req.Dir is the root repo (git-init'd
// by Setup with an initial commit); req.Path is "ext". Leaf Setup either wires
// ext/ in as a real submodule or leaves it a nested separate repo.
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	ok, err := IsSubmodule(req.Dir, req.Path)
	if err != nil {
		return nil, err
	}
	return &Response{Submodule: ok}, nil
}

var _ = os.MkdirTemp
var _ = exec.Command
var _ = filepath.Join
```
