# Scenario

**Feature**: modified tracked file dirties both clean APIs

```
# modify README.md after commit
IsDiffClean=false, IsPorcelainClean=false
```

## Preconditions

- Clean repo with tracked `README.md`.

## Steps

1. Overwrite `README.md` without staging/committing.

## Context

- Contrast leaf ensuring both APIs report dirty for tracked edits.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	writeFile(t, req.Dir, "README.md", "modified after commit\n")
	return nil
}
```
