# ListIgnoredDirs Test Case Tree

Run with:
```sh
doctest test ./ -v
```

## Version

0.0.2

## DSN (Domain Specific Notion)

The **caller** hands `ListIgnoredDirs(dir)` a directory inside a git work tree. The
**helper** runs `git -C <dir> ls-files --others --ignored --exclude-standard
--directory -z -- .` and returns the slash-relative paths of directories git reports
as ignored (no trailing slash, no `./`). When `dir` is not inside a git repo (or git
is unavailable) the helper returns an empty slice and a nil error — git-based skips
are simply inapplicable.

## Decision Tree

The most significant factor is **whether dir is a git repo** (controls whether git is
consulted at all); under that, the only branch is whether an ignored directory is
present and listed.

```
ListIgnoredDirs(dir)
├── ignored-dir/            # git repo with .gitignore=ignored/ -> ["ignored"]
└── not-a-git-repo/         # dir not a git repo -> []  (no error)
```

## How to Run

```sh
doctest vet ./gitops-dot/git/tests/list_ignored_dirs
doctest test ./gitops-dot/git/tests/list_ignored_dirs
```

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

type Request struct {
	Dir string
}

type Response struct {
	Dirs []string
}

// Run calls ListIgnoredDirs(req.Dir). req.Dir is set up by Setup (a git repo for
// the ignored-dir leaf, a plain non-repo dir for the not-a-git-repo leaf).
func Run(t *testing.T, req *Request) (*Response, error) {
	dirs, err := ListIgnoredDirs(req.Dir)
	if err != nil {
		return nil, err
	}
	return &Response{Dirs: dirs}, nil
}

var _ = filepath.Join
var _ = exec.Command
```
