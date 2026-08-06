# Scenario

**Feature**: nested directory resolves to the outer repository root without Git local env

```
# request starts below the repository root
caller -> ShowToplevel(<repo>/go-pkgs)

# returned output is the outer root, not the nested directory
caller <- <repo>\n
```

## Preconditions

- The request directory is `<repo>/go-pkgs`.
- No inherited `GIT_DIR` is present.

## Steps

1. Use the nested directory prepared by the parent setup.

## Context

- This leaf protects the existing successful baseline behavior.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.NestedDir == "" {
		t.Fatal("parent setup did not provide nested dir")
	}
	return nil
}
```
