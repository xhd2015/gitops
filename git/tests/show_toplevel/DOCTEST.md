# ShowToplevel Test Case Tree

## Version

0.0.2

# DSN (Domain Specific Notion)

`ShowToplevel(dir)` asks Git for the root worktree that contains `dir`.

A caller may run in a normal shell, where Git local environment variables are
absent. A caller may also run inside a Git hook, where variables such as
`GIT_DIR` are inherited from the hook runner. Those inherited variables are
ambient context: they must not change which repository root belongs to `dir`.

The helper preserves the raw command-style output contract and returns the
toplevel path with Git's trailing newline.

## Decision Tree

```text
show_toplevel
`-- git-local-env
    |-- absent
    |   `-- nested-dir-returns-root
    `-- inherited-git-dir
        `-- nested-dir-still-returns-root
```

## Test Case Index

| # | Path | Git local env | Expected |
|---|---|---|---|
| 1 | `git-local-env/absent/nested-dir-returns-root/` | unset | returns the outer repo root with trailing newline |
| 2 | `git-local-env/inherited-git-dir/nested-dir-still-returns-root/` | `GIT_DIR=<repo>/.git` | returns the outer repo root with trailing newline |

## How to Run

```sh
doctest vet git/tests/show_toplevel
doctest test git/tests/show_toplevel/...
```

```go
import (
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/xhd2015/doctest/session"
	gitpkg "github.com/xhd2015/gitops/git"
	"os/exec"
)

var showToplevelEnvMu sync.Mutex

type Request struct {
	RepoDir       string
	NestedDir     string
	ExpectedTop   string
	InheritGitDir bool
	GitDir        string
}

type Response struct {
	TopOutput string
	RawOutput string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	var raw string
	if req.InheritGitDir {
		out, err := exec.Command("git", "-C", req.NestedDir, "rev-parse", "--show-toplevel").Output()
		if err != nil {
			t.Fatalf("raw git proof without inherited env unexpectedly failed: %v", err)
		}
		_ = out

		cmd := exec.Command("git", "-C", req.NestedDir, "rev-parse", "--show-toplevel")
		cmd.Env = append(os.Environ(), "GIT_DIR="+req.GitDir)
		rawBytes, rawErr := cmd.Output()
		if rawErr != nil {
			t.Fatalf("raw git proof with inherited GIT_DIR unexpectedly failed: %v", rawErr)
		}
		raw = string(rawBytes)
	}

	showToplevelEnvMu.Lock()
	defer showToplevelEnvMu.Unlock()

	oldGitDir, hadGitDir := os.LookupEnv("GIT_DIR")
	if req.InheritGitDir {
		if err := os.Setenv("GIT_DIR", req.GitDir); err != nil {
			return nil, err
		}
	} else {
		if err := os.Unsetenv("GIT_DIR"); err != nil {
			return nil, err
		}
	}
	defer func() {
		if hadGitDir {
			_ = os.Setenv("GIT_DIR", oldGitDir)
		} else {
			_ = os.Unsetenv("GIT_DIR")
		}
	}()

	top, err := gitpkg.ShowToplevel(req.NestedDir)
	if err != nil {
		return nil, err
	}
	return &Response{TopOutput: top, RawOutput: strings.TrimSpace(raw)}, nil
}
```
