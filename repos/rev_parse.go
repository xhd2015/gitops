package repos

import (
	"strings"

	"github.com/xhd2015/gitops/git"
)

// RevParseVerified resolves a git ref to a commit hash.
// For bare-cache worktrees, origin/<branch> may be missing until fetch populates
// refs/remotes/origin/*; fall back to refs/heads/<branch> from the bare clone.
func RevParseVerified(dir, ref string) (string, error) {
	commit, err := git.RevParseVerified(dir, ref)
	if err == nil {
		return commit, nil
	}
	if !strings.HasPrefix(ref, "origin/") {
		return "", err
	}
	fallback := "refs/heads/" + strings.TrimPrefix(ref, "origin/")
	commit, fallbackErr := git.RevParseVerified(dir, fallback)
	if fallbackErr == nil {
		return commit, nil
	}
	return "", err
}