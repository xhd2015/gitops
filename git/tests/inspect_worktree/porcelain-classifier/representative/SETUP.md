# Scenario

**Feature**: classifyPorcelainLine maps porcelain XY pairs to change kinds

```
# table of XY inputs -> added/changed/renamed/deleted/none
TestExported_ClassifyPorcelainLine(xy) -> kind string
```

## Preconditions

- None; pure classifier function.

## Steps

1. Populate `req.PorcelainCases` with representative XY pairs.

## Context

- Rules: `??`→added; R→renamed; D→deleted; A→added; M→changed; else none.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.PorcelainCases = []PorcelainCase{
		{"??", "added"},
		{" M", "changed"},
		{"M ", "changed"},
		{"MM", "changed"},
		{" D", "deleted"},
		{"D ", "deleted"},
		{"R ", "renamed"},
		{" R", "renamed"},
		{"A ", "added"},
		{" A", "added"},
		{"AD", "deleted"},
		{"  ", "none"},
	}
	return nil
}
```