# Scenario

**Feature**: one staged file is returned correctly

```
# stage file.txt -> GetStagedFiles -> ["file.txt"]
write file.txt -> git add file.txt -> git diff --cached -> ["file.txt"]
```

## Preconditions

- A git repository with an initial commit.

## Steps

1. Create and stage `file.txt`.
2. Run `GetStagedFiles(req.Dir)`.

```go
func Setup(t *testing.T, req *Request) error {
	return stageFile(req.Dir, "file.txt", "hello")
}
```
