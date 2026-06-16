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
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xhd2015/gitops/git"
)

type Request struct {
	Dir  string
	RefA string
	RefB string
}

type Response struct {
	Result *git.CompareBranchesResult
}

func Setup(t *testing.T, req *Request) error {
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
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
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

func Run(t *testing.T, req *Request) (*Response, error) {
	result, err := git.CompareBranches(req.Dir, req.RefA, req.RefB)
	if err != nil {
		return nil, err
	}
	return &Response{Result: result}, nil
}
```
