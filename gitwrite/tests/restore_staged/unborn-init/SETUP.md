# Scenario

**Feature**: RestoreStaged unstages on unborn HEAD (`git init`, no commits)

```
# no HEAD to restore from: rm --cached is unstage
git init -> add drop.bin keep.txt -> RestoreStaged(drop.bin)
  -> drop.bin unstaged, still on disk; keep.txt still staged
```

## Preconditions

- Isolated repo after `git init` only (overwrites parent seed-commit dir).

## Steps

1. Init unborn repo (no commit).
2. Stage `keep.txt` and `drop.bin`.
3. Set `req.Paths = ["drop.bin"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	dir, err := initUnbornRepo(t)
	if err != nil {
		return err
	}
	req.Dir = dir
	if err := stageFile(req.Dir, "keep.txt", "keep\n"); err != nil {
		return err
	}
	if err := stageFile(req.Dir, "drop.bin", "\x00\x01\x02\x03"); err != nil {
		return err
	}
	req.Paths = []string{"drop.bin"}
	return nil
}
```
