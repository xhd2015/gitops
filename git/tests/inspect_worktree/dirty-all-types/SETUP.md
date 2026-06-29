# Scenario

**Feature**: InspectWorktree counts all porcelain change types

```
# baseline commit -> untracked + modified + git mv rename + delete
InspectWorktree -> dirty (1 added, 1 changed, 1 renamed, 1 deleted)
```

## Preconditions

- Clean repo from root setup.

## Steps

1. Add tracked files and commit baseline.
2. Create untracked file, modify tracked file, rename via `git mv`, delete tracked file.

## Context

- Untracked (`??`) counts as **added**.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/xgo/support/cmd"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir

	writeFile(t, dir, "tracked.txt", "baseline")
	writeFile(t, dir, "to-delete.txt", "baseline")
	writeFile(t, dir, "to-rename.txt", "baseline")
	gitAddA(t, dir)
	gitCommit(t, dir, "baseline")

	writeFile(t, dir, "untracked.txt", "new")
	writeFile(t, dir, "tracked.txt", "modified")
	if err := os.Remove(filepath.Join(dir, "to-delete.txt")); err != nil {
		return err
	}
	if err := cmd.Dir(dir).Run("git", "mv", "to-rename.txt", "renamed.txt"); err != nil {
		return err
	}
	return nil
}
```