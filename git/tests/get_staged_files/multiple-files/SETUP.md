# Scenario

**Feature**: multiple staged files are all returned

```
# stage a.txt b.txt -> GetStagedFiles -> ["a.txt", "b.txt"]
write a.txt b.txt -> git add a.txt b.txt -> git diff --cached -> ["a.txt", "b.txt"]
```

## Preconditions

- A git repository with an initial commit.

## Steps

1. Create and stage `a.txt` and `b.txt`.
2. Run `GetStagedFiles(req.Dir)`.

```go
func Setup(t *testing.T, req *Request) error {
	if err := stageFile(req.Dir, "a.txt", "aaa"); err != nil {
		return err
	}
	return stageFile(req.Dir, "b.txt", "bbb")
}
```
