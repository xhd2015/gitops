# Scenario

**Feature**: show toplevel works when Git local environment variables are absent

```
# normal shell does not provide GIT_DIR
shell without GIT_DIR -> ShowToplevel(<repo>/go-pkgs)

# helper returns the containing repository root
ShowToplevel <- <repo>\n
```

## Preconditions

- `GIT_DIR` is not part of the simulated caller environment.

## Steps

1. Ensure the request runs without inherited `GIT_DIR`.

## Context

- This is the baseline behavior current callers already rely on.

```go
func Setup(t *testing.T, req *Request) error {
	req.InheritGitDir = false
	return nil
}
```
