# Scenario

**Feature**: deleted staged files are excluded from result

```
# commit file, rm and stage deletion -> GetStagedFiles -> [] (excluded by --diff-filter=ACMRT)
git commit file.txt -> git rm file.txt -> git diff --cached --diff-filter=ACMRT -> empty
```

## Preconditions

- A git repository with an initial commit and `remove-me.txt` committed.

## Steps

1. Create and commit `remove-me.txt`.
2. Delete and stage the deletion with `git rm`.
3. Run `GetStagedFiles(req.Dir)`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"os/exec"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir
	// Create and commit the file
	if err := stageFile(dir, "remove-me.txt", "to be deleted"); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "add remove-me").Run(); err != nil {
		return err
	}
	// Delete and stage the deletion
	return exec.Command("git", "-C", dir, "rm", "remove-me.txt").Run()
}
```
