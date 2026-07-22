# Scenario

**Feature**: Unified for a rename with content change keeps rename meta plus
---/+++ and @@ hunks (non-empty Hunks)

```
# parse rename+edit raw → CachedDiff with Kind=rename and Hunks non-empty
ParseCachedDiff(rename+hunks raw)
  -> FilePatch{Kind:rename, OldPath, NewPath, Hunks:[...]}
  -> Unified() includes:
       diff --git a/old.go b/new.go
       rename from old.go
       rename to new.go
       --- a/old.go
       +++ b/new.go
       @@ ...
       hunk body
```

## Preconditions

- Mode `parse` — deterministic; live `git mv` + edit often becomes delete+add
  without stable rename detection.
- Raw is a valid unified diff with rename from/to and at least one `@@` hunk.

## Steps

1. Set `req.Mode = "parse"`.
2. Set `req.Raw` to a rename-with-content unified diff (old.go → new.go + hunk).
3. Run `git.ParseCachedDiff(req.Raw)`.
4. Assert seals `Unified()` still includes rename meta, ---/+++, and @@ body.

## Context

- Complements `pure-rename-unified/`: empty Hunks → no ---/+++; non-empty → keep
  content section (today's render already does this path).
- Similarity-index lines may be ignored by the parser; not required in Unified.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "parse"
	// Rename with content: Kind=rename and non-empty Hunks after parse.
	req.Raw = "" +
		"diff --git a/old.go b/new.go\n" +
		"similarity index 80%\n" +
		"rename from old.go\n" +
		"rename to new.go\n" +
		"index 1234567..abcdefg 100644\n" +
		"--- a/old.go\n" +
		"+++ b/new.go\n" +
		"@@ -1,2 +1,3 @@\n" +
		" context\n" +
		"-old line\n" +
		"+new line\n" +
		"+extra\n"
	return nil
}
```
