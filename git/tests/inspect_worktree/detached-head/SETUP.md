# Scenario

**Feature**: InspectWorktree when HEAD is detached

```
# git checkout --detach HEAD
InspectWorktree -> Branch="(detached)", commit still shown, IsClean=true
```

## Preconditions

- Repo with at least one commit.

## Steps

1. Run `git checkout --detach HEAD`.

## Context

- Commit hash and message must still be populated.

```go
import (
	"testing"

	"github.com/xhd2015/xgo/support/cmd"
)

func Setup(t *testing.T, req *Request) error {
	if err := cmd.Dir(req.Dir).Run("git", "checkout", "--detach", "HEAD"); err != nil {
		return err
	}
	return nil
}
```