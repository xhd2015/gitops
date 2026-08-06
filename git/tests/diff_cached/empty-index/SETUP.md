# Scenario

**Feature**: empty staging area returns (nil, nil)

```
# fresh repo, nothing staged -> DiffCached -> (nil, nil)
fresh repo (initial commit only) -> git diff --cached empty -> (nil, nil)
```

## Preconditions

- A git repository with an initial commit and no additional staged files.

## Steps

1. No additional files staged — just the initial commit from root SETUP.
2. Confirm the staging area is empty.
3. Run `git.DiffCached(req.Dir)`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"os/exec"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Root SETUP already created the repo with an initial commit.
	// Verify staging area is empty before the test runs.
	out, err := exec.Command("git", "-C", req.Dir, "diff", "--cached").Output()
	if err != nil {
		return err
	}
	if len(out) > 0 {
		t.Fatalf("precondition: staging area is not empty: %s", string(out))
	}
	return nil
}
```
