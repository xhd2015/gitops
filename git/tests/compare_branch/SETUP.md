## Preconditions
- `git` is available in PATH

## Steps
1. Create a temporary git repository with an initial commit on `main`
2. The repo directory is set as `req.Dir` for downstream setup steps

## Context
- Go module: `github.com/xhd2015/gitops`
- Package under test: `git`
- `CompareBranches(dir, refA, refB)` returns `(*CompareBranchesResult, error)`
- `CompareBranchesResult` has fields: `Relation`, `CommitsAheadA`, `CommitsAheadB`, `MergeBase`, `DiffFileCount`
- `BranchRelation` values: `BranchRelationSame`, `BranchRelationAIsAncestorOfB`, `BranchRelationBIsAncestorOfA`, `BranchRelationDiverged`

```go
import (
	"os"
	"testing"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
	"github.com/xhd2015/gitops/git"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir, err := os.MkdirTemp("", "gitops-compare-branch-*")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	runGit := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v failed: %v
%s", args, err, out)
		}
	}

	runGit("init")
	runGit("config", "user.email", "test@test.com")
	runGit("config", "user.name", "test")

	err = os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test"), 0644)
	if err != nil {
		return err
	}
	runGit("add", ".")
	runGit("commit", "-m", "initial commit")
	runGit("branch", "-M", "main")

	req.Dir = dir
	return nil
}
```
