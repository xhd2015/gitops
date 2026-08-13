package repos

import (
	"context"
	"fmt"
	"strings"
)

// LsRemoteBranchTip resolves the tip commit SHA of a remote branch via
// `git ls-remote` (advertise only — no object pack). branch is a simple name
// like "master" or "feat/x" (no refs/heads/ prefix required).
func LsRemoteBranchTip(ctx context.Context, cloneURL, branch string, auth *GitAuthConfig, env map[string]string) (string, error) {
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return "", fmt.Errorf("branch is required")
	}
	if strings.HasPrefix(branch, "refs/") {
		return "", fmt.Errorf("branch must be a short name, not %q", branch)
	}
	cloneURL = strings.TrimSpace(cloneURL)
	if cloneURL == "" {
		return "", fmt.Errorf("clone URL is required")
	}

	merged, gitURL, err := mergeAuth(auth, cloneURL)
	if err != nil {
		return "", fmt.Errorf("parse auth: %w", err)
	}
	gitURL = cloneURLForGit(gitURL, merged)

	ref := "refs/heads/" + branch
	out, err := gitOutputWithEnvContext(ctx, "", merged, env, "ls-remote", gitURL, ref)
	if err != nil {
		return "", fmt.Errorf("ls-remote %s: %s", branch, merged.maskError(err))
	}

	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		// fields[0]=sha fields[1]=ref (may be refs/heads/X^{} peeled — prefer non-peeled)
		name := fields[1]
		if strings.HasSuffix(name, "^{}") {
			continue
		}
		if name == ref || strings.HasSuffix(name, "/"+branch) {
			sha := fields[0]
			if sha == "" {
				continue
			}
			return sha, nil
		}
	}
	return "", fmt.Errorf("branch %q not found on remote", branch)
}
