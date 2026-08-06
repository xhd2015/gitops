## Expected

- `DiffCached` returns a non-nil `*model.CachedDiff` with no error.
- Full `Unified()` contains `long.txt` and is longer than the truncated form.
- `UnifiedTruncated(24)`:
  - contains the omit marker pattern `more lines omitted`
  - matches commit_msg style `...(` + digits + ` more lines omitted)`
  - does **not** include the late body marker `LATE_MARKER_LINE_39` when that
    marker is present in full `Unified()`
  - has strictly fewer lines than `Unified()`

```go
import (
	"regexp"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error for long staged file: %v", err)
	}
	if resp == nil || resp.Diff == nil {
		t.Fatal("expected non-nil Diff for long staged file")
	}

	full := resp.Diff.Unified()
	trunc := resp.Diff.UnifiedTruncated(24)

	if full == "" {
		t.Fatal("Unified() returned empty string")
	}
	if !strings.Contains(full, "long.txt") {
		t.Fatalf("Unified() missing path long.txt; got:\n%s", full)
	}

	if trunc == "" {
		t.Fatal("UnifiedTruncated(24) returned empty string")
	}
	if !strings.Contains(trunc, "more lines omitted") {
		t.Fatalf("UnifiedTruncated(24) missing omit marker; got:\n%s", trunc)
	}
	// Policy: "...(N more lines omitted)" with N = remaining lines.
	omitRE := regexp.MustCompile(`\.\.\.\(\d+ more lines omitted\)`)
	if !omitRE.MatchString(trunc) {
		t.Fatalf("UnifiedTruncated(24) omit marker shape mismatch; got:\n%s", trunc)
	}

	fullLines := strings.Split(strings.TrimRight(full, "\n"), "\n")
	truncLines := strings.Split(strings.TrimRight(trunc, "\n"), "\n")
	if len(truncLines) >= len(fullLines) {
		t.Fatalf("expected truncated line count < full; trunc=%d full=%d",
			len(truncLines), len(fullLines))
	}

	const late = "LATE_MARKER_LINE_39"
	if strings.Contains(full, late) && strings.Contains(trunc, late) {
		t.Fatalf("UnifiedTruncated(24) still contains late body %q; got:\n%s", late, trunc)
	}
	// Precondition soft-check: full unified should include the late marker so
	// the omit assertion is meaningful (staging produced enough body lines).
	if !strings.Contains(full, late) {
		t.Fatalf("precondition: full Unified should contain %q so truncation is observable; full:\n%s", late, full)
	}
}
```
