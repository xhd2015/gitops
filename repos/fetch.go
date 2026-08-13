package repos

import (
	"context"
	"fmt"
	"strings"
)

// URLFetchOptions configures FetchFromURL.
type URLFetchOptions struct {
	Refs      []string          // SHAs / refs to fetch
	Depth     int               // 0 = no --depth; >0 = --depth=N
	Deepen    int               // >0 = --deepen=N (history growth on shallow bare)
	Unshallow bool              // true = --unshallow when repo is still shallow
	Env       map[string]string // MergeGitEnv into subprocess
}

// FetchFromURL fetches refs into dir using cloneURL as the remote argument (URL-as-is).
// Credentials ride in the URL userinfo when present; this path does not use
// http.extraHeader / GitAuthConfig.gitArgs. Tokens are masked in returned errors.
//
// Depth, Deepen, and Unshallow are mutually preferred as: Unshallow > Deepen > Depth.
//
// This is an unlocked primitive. For shared static-cache bares, call only while
// holding the repo flock (see FetchCommits / EnsureStaticBareCommits) so concurrent
// shallow fetches cannot race on git's shallow.lock.
func FetchFromURL(ctx context.Context, dir, cloneURL string, opts URLFetchOptions) error {
	if cloneURL == "" {
		return fmt.Errorf("fetch from url: cloneURL is required")
	}

	args := []string{"fetch"}
	switch {
	case opts.Unshallow:
		shallow, _ := gitOutput(dir, "rev-parse", "--is-shallow-repository")
		if strings.TrimSpace(shallow) == "true" {
			args = append(args, "--unshallow")
		}
	case opts.Deepen > 0:
		args = append(args, fmt.Sprintf("--deepen=%d", opts.Deepen))
	case opts.Depth > 0:
		args = append(args, fmt.Sprintf("--depth=%d", opts.Depth))
	}
	// Remote is the clone URL as-is (not bare name "origin").
	args = append(args, cloneURL)
	for _, ref := range opts.Refs {
		if ref == "" {
			continue
		}
		args = append(args, ref)
	}

	// auth=nil: no -c http.extraHeader; env merged via runGitWithEnvContext.
	if err := runGitWithEnvContext(ctx, dir, nil, opts.Env, args...); err != nil {
		return fmt.Errorf("fetch from url: %s", maskFetchURLError(err, cloneURL))
	}
	return nil
}

// maskFetchURLError redacts token material from cloneURL userinfo in error text.
func maskFetchURLError(err error, cloneURL string) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	info, extractErr := extractTokenFromURL(cloneURL)
	if extractErr != nil || info == nil || !info.HasToken || info.Token == "" {
		return msg
	}
	return maskSensitiveInfo(msg, info.Token)
}
