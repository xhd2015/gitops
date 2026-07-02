# Scenario

**Feature**: no staged files returns an empty slice

```
# fresh repo, nothing staged -> GetStagedFiles -> []
fresh repo (initial commit only) -> git diff --cached = empty -> []
```

## Preconditions

- A git repository with an initial commit and no additional staged files.

## Steps

1. No additional files staged — just the initial commit.
2. Run `GetStagedFiles(req.Dir)`.

```go
import "os/exec"

func Setup(t *testing.T, req *Request) error {
	// Root SETUP already created the repo with an initial commit.
	// Verify staging area is empty before the test runs.
	out, err := exec.Command("git", "-C", req.Dir, "diff", "--cached", "--name-only").Output()
	if err != nil {
		return err
	}
	if len(out) > 0 {
		t.Logf("precondition: staging area is not empty: %s", string(out))
	}
	return nil
}
```
