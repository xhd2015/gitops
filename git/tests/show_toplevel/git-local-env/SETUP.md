# Scenario

**Feature**: show toplevel is stable across Git local environment modes

```
# caller may run with no local Git environment or inside a Git hook
caller environment -> ShowToplevel(<repo>/go-pkgs)

# helper identifies the repository from the requested directory, not ambient hook state
ShowToplevel -> git repository at <repo>
```

## Preconditions

- The parent setup created a repository with a nested directory.

## Steps

1. Select the Git local environment mode in a child scenario.

## Context

- This grouping separates normal shell execution from inherited hook-style Git environment.

```go
func Setup(t *testing.T, req *Request) error {
	req.InheritGitDir = false
	return nil
}
```
