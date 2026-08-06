# Scenario

**Feature**: Unified returns a full unified-diff string for the staged patch

```
# stage small note.txt -> DiffCached -> Unified() contains "diff --git" and path
stage note.txt -> DiffCached(dir) -> d.Unified() includes "diff --git" + "note.txt"
```

## Preconditions

- A git repository with an initial commit.
- One small text file will be staged.

## Steps

1. Stage `note.txt` with a short body.
2. Run `git.DiffCached(req.Dir)`.
3. Assert calls `resp.Diff.Unified()` and checks for unified-diff markers.

## Context

- Classic RED: `(*model.CachedDiff).Unified` does not exist yet.
- "Roundtrip-ish": we do not require byte-identical git stdout, only that the
  rendered text is a usable unified diff mentioning the staged path.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	return stageFile(req.Dir, "note.txt", "hello from note\n")
}
```
