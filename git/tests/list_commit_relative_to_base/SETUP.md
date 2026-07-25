# Scenario

**Feature**: branch scanning recovers only direct merge side chains.

```text
# a test constructs a Git commit DAG in an isolated repository
test graph -> ListCommitRelativeToBase -> returned branch commits
```

## Preconditions

- `git` is available on PATH.

## Steps

1. Create an isolated Git repository for the leaf.
2. Use `git commit-tree` to construct the leaf's commit DAG.

## Context

- Commit subjects are stable graph labels used by assertions.

```go
import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "testing"

    "github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    dir, err := os.MkdirTemp("", "gitops-list-relative-")
    if err != nil {
        return err
    }
    req.Dir = dir
    if _, err := gitOutput(dir, "init", "--quiet"); err != nil {
        return err
    }
    return nil
}

func gitOutput(dir string, args ...string) (string, error) {
    command := exec.Command("git", args...)
    command.Dir = dir
    output, err := command.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, output)
    }
    return strings.TrimSpace(string(output)), nil
}

func commit(t *testing.T, dir string, subject string, parents ...string) string {
    t.Helper()
    tree, err := gitOutput(dir, "mktree")
    if err != nil {
        t.Fatal(err)
    }
    args := []string{"commit-tree", tree, "-m", subject}
    for _, parent := range parents {
        args = append(args, "-p", parent)
    }
    hash, err := gitOutput(dir, args...)
    if err != nil {
        t.Fatal(err)
    }
    return hash
}
```
