# Scenario

**Feature**: ParseCachedDiff rejects non-empty unparseable raw with typed error

```
# pure parse of malformed unified-diff text (no live git stdout injection)
caller -> ParseCachedDiff("not a valid unified diff...")
  -> (nil, *DiffCachedParseError{Raw: full input})
```

## Preconditions

- Exported helper `git.ParseCachedDiff(raw string) (*model.CachedDiff, error)`
  is the entrypoint under test (same parse failure contract as DiffCached).
- `git.DiffCachedParseError` exists with `Raw` (and Error/Unwrap for As).

## Steps

1. Set `req.Mode` to `parse`.
2. Provide non-empty malformed `req.Raw` that is not a valid staged unified diff.
3. Run calls `git.ParseCachedDiff(req.Raw)`.

## Context

- Prefer ParseCachedDiff over injecting fake git stdout so this leaf is
  deterministic Classic RED until the parser exists.
- Live DiffCached should wrap the same parse path; this leaf seals the parse
  failure surface independently.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "parse"
	// Non-empty garbage: not a unified diff header/hunk stream.
	req.Raw = "not a valid unified diff\n@@ broken\nthis is not a patch\n"
	return nil
}
```
