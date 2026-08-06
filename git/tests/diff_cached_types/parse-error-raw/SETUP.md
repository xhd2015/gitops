# Scenario

**Feature**: git.DiffCachedParseError carries Dir/Raw and supports Error/Unwrap/As

```
# construct typed parse failure with raw staged patch text
DiffCachedParseError{Dir, Raw, Err: cause}
  -> Error() non-empty, mentions parse or diff
  -> Unwrap() == cause
  -> errors.As recovers Dir and Raw
```

## Preconditions

- Type `git.DiffCachedParseError` exists with fields `Dir`, `Raw`, `Err`.
- Methods: `Error() string`, `Unwrap() error` (or equivalent for `errors.As`/`Unwrap`).

## Steps

1. Set `req.Mode` to `parse-error-raw`.
2. Provide `Dir`, `Raw` (sample unparsable patch text), and `CauseMsg`.
3. Run constructs `*git.DiffCachedParseError` and records Error/Unwrap/As.

## Context

- Classic RED: type does not exist yet → build/runtime RED until implementer adds it.
- Raw is preserved for diagnostics even when parse fails.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "parse-error-raw"
	req.Dir = "/tmp/repo-under-test"
	req.Raw = "not a valid unified diff\n@@ broken"
	req.CauseMsg = "unexpected hunk header"
	return nil
}
```
