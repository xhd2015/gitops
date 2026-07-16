# Scenario

**Feature**: UnifiedTruncated caps each file section and appends an omit marker

```
# stage long file (>24 lines in unified section) -> UnifiedTruncated(24)
# keeps first 24 lines of the file section; appends "...(N more lines omitted)"
stage long.txt (many lines) -> DiffCached -> d.UnifiedTruncated(24)
  -> omit marker present; full Unified has more lines / late body
```

## Preconditions

- A git repository with an initial commit.
- Staged file body is long enough that the unified-diff *file section* exceeds
  24 lines (headers + hunk + body).

## Steps

1. Build a multi-line body (40+ content lines with a unique late marker).
2. Stage `long.txt` with that body.
3. Run `git.DiffCached(req.Dir)`.
4. Assert calls `Unified()` and `UnifiedTruncated(24)`.

## Context

- Classic RED: `(*model.CachedDiff).UnifiedTruncated` (and `Unified`) do not
  exist yet.
- Truncate policy matches old commit_msg: per file section, keep first
  `maxLinesPerFile` lines; if more, append `\n...(%d more lines omitted)`.
- Late unique line `LATE_MARKER_LINE_39` must appear in full Unified when
  present in the staged patch, but must be omitted from truncated output.

```go
import (
	"strings"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	var b strings.Builder
	for i := 1; i <= 40; i++ {
		if i == 39 {
			b.WriteString("LATE_MARKER_LINE_39\n")
			continue
		}
		b.WriteString(strings.Repeat("x", 8))
		b.WriteByte('\n')
	}
	return stageFile(req.Dir, "long.txt", b.String())
}
```
