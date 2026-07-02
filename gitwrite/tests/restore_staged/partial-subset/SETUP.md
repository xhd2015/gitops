# Scenario

**Feature**: unstaging a subset leaves other staged files untouched

```
# stage a.txt b.txt -> restore only a.txt -> a.txt gone, b.txt still staged
git add a.txt b.txt -> git restore --staged -- a.txt -> diff --cached = [b.txt]
```

## Preconditions

- A git repository with an initial commit.
- `a.txt` and `b.txt` have been created and staged.

## Steps

1. Create and stage `a.txt` and `b.txt`.
2. Set `req.Paths = ["a.txt"]` (only unstage `a.txt`).
3. Run `RestoreStaged`, then verify only `a.txt` is unstaged.

```go
func Setup(t *testing.T, req *Request) error {
	if err := stageFile(req.Dir, "a.txt", "aaa"); err != nil {
		return err
	}
	if err := stageFile(req.Dir, "b.txt", "bbb"); err != nil {
		return err
	}
	req.Paths = []string{"a.txt"}
	return nil
}
```
