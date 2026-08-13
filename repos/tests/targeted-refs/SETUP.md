# Scenario

**Feature**: targeted fetch dest refs use a customizable prefix (default `refs/gitops/targets`) and a conservative blob-less fetch into an injectable `ReposRoot`.

```
# per-leaf file:// remote + injectable ReposRoot (no HOME / Setenv / Chdir)
test harness -> seed master+feature under d.DOCTEST_CASE
             -> FetchTargetedRefs(url, refs, {ReposRoot, TargetRefPrefix, Progress})
             -> dest {prefix}/{sha256(input)} in targeted-cache

```

## Preconditions

- Nested doctest root under the gitops checkout (`repos/tests/targeted-refs`).
- Package under test: `github.com/xhd2015/gitops/repos`.
- `git` available on PATH (fixture uses `git_isolated`).
- Every leaf injects absolute `ReposRoot` and a local `file://` remote under `d.DOCTEST_CASE`.
- **Forbid** `t.Setenv` / `os.Setenv` / `os.Chdir` / `t.Chdir`.
- No network; no e2e label.
- Harness already parallels leaves; do not call `t.Parallel()` in Setup.

## Steps

1. Require `d.DOCTEST_CASE`; create `{DOCTEST_CASE}/.repos` as `ReposRoot`.
2. Seed a two-branch (`master`, `feature`) bare remote under `{DOCTEST_CASE}/fixture`.
3. Default request: fetch `master`, empty `TargetRefPrefix`, `Depth=1`.
4. Descendants narrow prefix / refs.

## Context

- Dest ref: `{normalizedPrefix}/{hex(sha256(input))}`.
- Default prefix: `refs/gitops/targets`.
- Cache: `{ReposRoot}/targeted-cache/…` (no second cache root).
- Root `Run` calls only `repos.FetchTargetedRefs` (then lists dest refs for Assert).

```go
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/gitops/git/git_isolated"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	if d == nil || d.DOCTEST_CASE == "" {
		return fmt.Errorf("missing d.DOCTEST_CASE")
	}
	req.ReposRoot = filepath.Join(d.DOCTEST_CASE, ".repos")
	if err := os.MkdirAll(req.ReposRoot, 0o755); err != nil {
		return err
	}
	remote, err := seedTwoBranchRemote(filepath.Join(d.DOCTEST_CASE, "fixture"))
	if err != nil {
		return err
	}
	req.CloneURL = "file://" + filepath.ToSlash(remote)
	req.Refs = []string{"master"}
	req.TargetRefPrefix = ""
	req.Depth = 1
	return nil
}

func seedTwoBranchRemote(fixtureRoot string) (string, error) {
	if err := os.MkdirAll(fixtureRoot, 0o755); err != nil {
		return "", err
	}
	work := filepath.Join(fixtureRoot, "work")
	remote := filepath.Join(fixtureRoot, "remote.git")
	if err := git_isolated.Init(work, "master"); err != nil {
		return "", err
	}
	if err := writeFile(filepath.Join(work, "README"), "master\n"); err != nil {
		return "", err
	}
	if err := git_isolated.Run(work, "add", "README"); err != nil {
		return "", err
	}
	if err := git_isolated.Run(work, "commit", "-m", "master"); err != nil {
		return "", err
	}
	if err := git_isolated.Run(work, "checkout", "-b", "feature"); err != nil {
		return "", err
	}
	if err := writeFile(filepath.Join(work, "feature.txt"), "feature\n"); err != nil {
		return "", err
	}
	if err := git_isolated.Run(work, "add", "feature.txt"); err != nil {
		return "", err
	}
	if err := git_isolated.Run(work, "commit", "-m", "feature"); err != nil {
		return "", err
	}
	if err := git_isolated.Run(work, "clone", "--bare", ".", remote); err != nil {
		return "", err
	}
	if err := git_isolated.Run(remote, "config", "uploadpack.allowFilter", "true"); err != nil {
		return "", err
	}
	if err := git_isolated.Run(remote, "config", "uploadpack.allowAnySHA1InWant", "true"); err != nil {
		return "", err
	}
	return remote, nil
}

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func expectedDestRef(prefix, input string) string {
	p := strings.TrimRight(prefix, "/")
	if p == "" {
		p = DefaultTargetRefPrefix
	}
	sum := sha256.Sum256([]byte(input))
	return p + "/" + hex.EncodeToString(sum[:])
}

func hasRef(refs []string, want string) bool {
	for _, r := range refs {
		if r == want {
			return true
		}
	}
	return false
}

func hasRefPrefix(refs []string, prefix string) bool {
	for _, r := range refs {
		if strings.HasPrefix(r, prefix) {
			return true
		}
	}
	return false
}

func destHashOf(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}
```
