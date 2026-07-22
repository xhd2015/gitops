# Scenario

**Feature**: Unified for a pure rename (Kind=rename, empty Hunks) emits only
diff --git + rename from/to — no ---/+++ file headers and no @@ hunks

```
# live 100% rename: git mv stages path-only change
git mv oldname.txt newname.txt
  -> DiffCached(dir) -> FilePatch{Kind:rename, Hunks:[]}
  -> Unified() =
       "diff --git a/oldname.txt b/newname.txt\n"
       "rename from oldname.txt\n"
       "rename to newname.txt"
  # must NOT include "--- a/" / "+++ b/" or "@@"
```

## Preconditions

- A git repository with an initial commit that includes `oldname.txt`.
- Staging a pure rename via `git mv` (similarity 100%, no content change).
- Parser maps rename meta to `Kind == "rename"` with empty `Hunks`.

## Steps

1. Write and commit is already done by root Setup (README only) — create and
   commit `oldname.txt` on top of the initial commit.
2. `git mv oldname.txt newname.txt` so the index holds a pure rename.
3. Run `git.DiffCached(req.Dir)`.
4. Assert calls `resp.Diff.Unified()` and seals pure-rename render rules.

## Context

- Classic RED until `renderFilePatch` skips `---/+++` when Kind=rename and
  `len(Hunks)==0` (and not Binary). Today Unified always emits `---/+++`
  after rename meta even for empty hunks.
- Live git is preferred: 100% renames omit ---/+++ in git stdout; re-render
  must match that shape for consumers (e.g. commit_msg).

```go
import (
	"os/exec"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	// Commit a tracked file, then pure-rename it in the index.
	if err := writeFile(req.Dir, "oldname.txt", "pure rename body\n"); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", req.Dir, "add", "oldname.txt").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", req.Dir, "commit", "-m", "add oldname").Run(); err != nil {
		return err
	}
	if err := exec.Command("git", "-C", req.Dir, "mv", "oldname.txt", "newname.txt").Run(); err != nil {
		return err
	}
	return nil
}
```
