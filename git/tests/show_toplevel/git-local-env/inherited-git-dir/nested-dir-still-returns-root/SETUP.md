# Scenario

**Bug**: nested directory must not become the toplevel just because GIT_DIR is inherited

```
# inherited Git hook env points at the outer repo gitdir
GIT_DIR=<repo>/.git -> ShowToplevel(<repo>/go-pkgs)

# safe helper reports the outer repo root
ShowToplevel <- <repo>\n
```

## Preconditions

- `GIT_DIR` is set to the temporary repository's `.git` directory.
- The request directory is the nested `<repo>/go-pkgs` directory.

## Steps

1. Use the hook-style `GIT_DIR` mode selected by the parent setup.

## Context

- This leaf is expected to be RED before implementation because the current helper inherits `GIT_DIR`.

```go
func Setup(t *testing.T, req *Request) error {
	if req.GitDir == "" {
		t.Fatal("parent setup did not provide git dir")
	}
	return nil
}
```
