# Scenario

**Feature**: after unstaging everything, DiffCached is (nil, nil) again

```
# stage file.txt -> restore --staged -> DiffCached -> (nil, nil)
git add file.txt -> git restore --staged -- . -> empty staged patch -> (nil, nil)
```

## Preconditions

- A git repository with an initial commit.

## Steps

1. Create and stage `file.txt`.
2. Confirm something is staged.
3. Unstage all with `git restore --staged -- .`.
4. Run `git.DiffCached(req.Dir)`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"os/exec"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if err := stageFile(req.Dir, "file.txt", "hello"); err != nil {
		return err
	}
	// Confirm something was staged, then unstage everything.
	out, err := exec.Command("git", "-C", req.Dir, "diff", "--cached", "--name-only").Output()
	if err != nil {
		return err
	}
	if len(out) == 0 {
		t.Fatal("precondition: expected file.txt to be staged before unstage")
	}
	if err := exec.Command("git", "-C", req.Dir, "restore", "--staged", "--", ".").Run(); err != nil {
		return err
	}
	return nil
}
```
