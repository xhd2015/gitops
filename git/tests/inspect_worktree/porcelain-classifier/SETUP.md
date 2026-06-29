# Scenario

**Feature**: porcelain XY status classification grouping

```
# leaves under porcelain-classifier/ exercise classifyPorcelainLine rules
```

## Preconditions

- Classifier tested via `TestExported_ClassifyPorcelainLine`.

## Steps

1. Descendant leaves supply XY pairs and expected kinds.

## Context

- Unexported `classifyPorcelainLine` exposed for doctests only.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Log("porcelain classifier leaves supply XY cases in leaf Setup")
	return nil
}
```