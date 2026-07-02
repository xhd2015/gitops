# Scenario

**Bug**: inherited GIT_DIR from a Git hook can corrupt raw `git -C dir rev-parse --show-toplevel`

```
# Git hook exports GIT_DIR for the repository being committed
git hook environment -> GIT_DIR=<repo>/.git

# caller asks about a nested directory in the same repository
caller -> ShowToplevel(<repo>/go-pkgs)

# helper must ignore ambient Git local env while preserving raw output shape
caller <- <repo>\n
```

## Preconditions

- The parent setup created `<repo>` and `<repo>/go-pkgs`.

## Steps

1. Simulate inherited hook state with `GIT_DIR=<repo>/.git`.
2. Call `ShowToplevel(<repo>/go-pkgs)` under that inherited environment.
3. Capture raw Git output under the same inherited `GIT_DIR` to prove the Git behavior that the helper must guard against.

## Context

- Raw Git returns `<repo>/go-pkgs` for this fixture when `GIT_DIR=<repo>/.git` is inherited.
- The safe helper should still return `<repo>\n`.

```go
func Setup(t *testing.T, req *Request) error {
	req.InheritGitDir = true
	return nil
}
```
