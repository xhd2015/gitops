## Preconditions
- A git repository exists at `dir`
- `git` is available in PATH

## Steps
1. Create a temporary directory
2. Run `git init` and `git branch -M master`
3. Write `README.md` with content "test" and commit as initial commit

## Context
- Go module: `github.com/xhd2015/gitops`
- Package under test: `git`
- The `Request.Dir` field is set to the temp repo path automatically by Setup

```go
type Request struct {
	Dir          string
	CompareWith  string // empty string means no compare
	ResolvePaths bool
}

type Response struct {
	Files []string
}
```
