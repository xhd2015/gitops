# Scenario

**Feature**: file in a subdirectory is returned with the correct relative path

```
# stage sub/pkg.go -> GetStagedFiles -> ["sub/pkg.go"]
write sub/pkg.go -> git add sub/pkg.go -> git diff --cached -> ["sub/pkg.go"]
```

## Preconditions

- A git repository with an initial commit.

## Steps

1. Create and stage `sub/pkg.go` inside a subdirectory.
2. Run `GetStagedFiles(req.Dir)`.

```go
func Setup(t *testing.T, req *Request) error {
	return stageFile(req.Dir, "sub/pkg.go", "package pkg")
}
```
