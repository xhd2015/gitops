# Scenario

**Feature**: DiffCached returns a structured staged patch (`*model.CachedDiff`)
with optional FileCount / Unified / UnifiedTruncated helpers

```
# caller passes a repo dir; DiffCached runs git and parses staged patch
caller -> DiffCached(dir) -> git -C dir diff --cached -> parse -> *CachedDiff

# empty staging area yields nil pointer, not an empty struct and not an error
git index empty vs HEAD -> stdout "" -> (nil, nil)

# staged files yield Files with Kind and paths
staged file.txt -> non-nil CachedDiff.Files

# non-git directory yields ordinary error
plain dir (no .git) -> git fails -> (nil, err)

# pure parse path for malformed raw (no live git injection)
caller -> ParseCachedDiff(raw) -> *DiffCachedParseError{Raw}

# render helpers on *CachedDiff (P3)
*CachedDiff -> FileCount() / Unified() / UnifiedTruncated(n)
nil *CachedDiff -> FileCount 0, Unified "", UnifiedTruncated ""

# pure rename (Kind=rename, empty Hunks): rename meta only
FilePatch{Kind:rename, Hunks:[]} -> Unified = diff --git + rename from/to
  (no ---/+++, no @@)

# rename with content (Kind=rename, Hunks non-empty): meta + content section
FilePatch{Kind:rename, Hunks:[...]} -> Unified keeps ---/+++ + @@ hunks
```

## Preconditions

- `git` is available in PATH for live leaves.
- Default leaves use a git repository with an initial commit (staged diffs are
  relative to HEAD). The `not-a-git-repo` leaf overrides `req.Dir`.
- The `parse-malformed` leaf sets `req.Mode = "parse"` and does not rely on
  the temp repo for assertions.
- The `nil-receiver-methods` leaf sets `req.Mode = "nil-methods"` and does not
  call git.
- Locked contract: `func DiffCached(dir string) (*model.CachedDiff, error)`.
- Parse helper (for testability): `func ParseCachedDiff(raw string) (*model.CachedDiff, error)`.
- P3 methods on `*model.CachedDiff`:
  - `func (d *CachedDiff) FileCount() int`
  - `func (d *CachedDiff) Unified() string`
  - `func (d *CachedDiff) UnifiedTruncated(maxLinesPerFile int) string`
- P1 pure-rename render (classic RED until `renderFilePatch` special-cases
  empty-hunk renames): no `---/+++` / no `@@` when Kind=rename, !Binary,
  len(Hunks)==0; rename+hunks keeps content section.
- Out of scope: commit_msg integration, similarity-index storage.

## Steps

1. Create a temporary directory and git-init it with an initial commit.
2. Each leaf stages, unstages, replaces `req.Dir`, switches to parse mode, or
   switches to nil-methods mode.
3. Run calls `git.DiffCached(req.Dir)`, `git.ParseCachedDiff(req.Raw)`, or
   returns a nil Diff for nil-methods.

## Context

- Go module: `github.com/xhd2015/gitops`
- Package under test: `git` and `model`
- Model: `model.CachedDiff` / `model.FilePatch` / `model.Hunk`
- Typed parse failure: `git.DiffCachedParseError` (Dir, Raw, Err)
- Truncate policy (match old commit_msg): per file section keep first
  `maxLinesPerFile` lines; if more, append `\n...(%d more lines omitted)`.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir, err := os.MkdirTemp("", "gittest-diff-cached")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	if err := exec.Command("git", "-C", dir, "init").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "branch", "-M", "master").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "config", "user.email", "test@example.com").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "config", "user.name", "Test User").Run(); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("test"), 0644); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "init").Run(); err != nil {
		return err
	}
	req.Dir = dir
	return nil
}

func writeFile(dir string, name string, content string) error {
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func stageFile(dir string, name string, content string) error {
	if err := writeFile(dir, name, content); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "add", name).Run()
}
```
