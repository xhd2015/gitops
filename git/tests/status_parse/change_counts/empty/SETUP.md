# Scenario

**Feature**: empty porcelain yields zero ChangeCounts

```
"" -> ParseChangeCounts -> all zeros
```

## Preconditions

- None.

## Steps

1. Set empty porcelain string.

## Context

- Clean baseline for four-bucket parser.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "change_counts"
	req.Porcelain = ""
	return nil
}
```
