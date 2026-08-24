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

	// NoBlobFilter disables --filter=blob:none on Deepen/Unshallow. Default
	// (false) applies the filter for history-growth fetches so connectivity
	// ladders stay small; Depth fetches never use the filter (worktrees need
	// blobs for the fetched tip). EnsureStaticBareCommits sets this on its
	// one-shot unshallow so worktree/cat-file do not dial a promisor remote.
	NoBlobFilter bool
}

// FetchFromURL fetches refs into dir using cloneURL as the remote argument (URL-as-is).
// Credentials ride in the URL userinfo when present; this path does not use
// http.extraHeader / GitAuthConfig.gitArgs. Tokens are masked in returned errors.
//
// Depth, Deepen, and Unshallow are mutually preferred as: Unshallow > Deepen > Depth.
//
// Deepen and Unshallow default to --filter=blob:none (partial history for
// merge-base / connectivity). Depth-N tip fetches stay unfiltered so detached
// worktrees can materialize trees. If the remote rejects filtering, Deepen/
// Unshallow retry once without the filter.
//
// This is an unlocked primitive. For shared static-cache bares, call only while
// holding the repo flock (see FetchCommits / EnsureStaticBareCommits) so concurrent
// shallow fetches cannot race on git's shallow.lock.
func FetchFromURL(ctx context.Context, dir, cloneURL string, opts URLFetchOptions) error {
	if cloneURL == "" {
		return fmt.Errorf("fetch from url: cloneURL is required")
	}

	wantFilter := !opts.NoBlobFilter && (opts.Deepen > 0 || opts.Unshallow)
	args := buildURLFetchArgs(dir, cloneURL, opts, wantFilter)
	if len(args) == 0 {
		// Unshallow requested on a non-shallow repo: nothing to fetch.
		return nil
	}

	// auth=nil: no -c http.extraHeader; env merged via runGitWithEnvContext.
	if err := runGitWithEnvContext(ctx, dir, nil, opts.Env, args...); err != nil {
		if wantFilter && filterRejected(err) {
			retry := buildURLFetchArgs(dir, cloneURL, opts, false)
			if len(retry) == 0 {
				return nil
			}
			if retryErr := runGitWithEnvContext(ctx, dir, nil, opts.Env, retry...); retryErr != nil {
				return fmt.Errorf("fetch from url: %s", maskFetchURLError(retryErr, cloneURL))
			}
			return nil
		}
		return fmt.Errorf("fetch from url: %s", maskFetchURLError(err, cloneURL))
	}
	return nil
}

// buildURLFetchArgs returns git fetch argv (without "git"). Empty means no-op
// (unshallow on an already-full history).
func buildURLFetchArgs(dir, cloneURL string, opts URLFetchOptions, withBlobNone bool) []string {
	args := []string{"fetch"}
	if withBlobNone {
		args = append(args, "--filter=blob:none")
	}
	switch {
	case opts.Unshallow:
		shallow, _ := gitOutput(dir, "rev-parse", "--is-shallow-repository")
		if strings.TrimSpace(shallow) != "true" {
			return nil
		}
		args = append(args, "--unshallow")
	case opts.Deepen > 0:
		args = append(args, fmt.Sprintf("--deepen=%d", opts.Deepen))
	case opts.Depth > 0:
		args = append(args, fmt.Sprintf("--depth=%d", opts.Depth))
	}
	args = append(args, cloneURL)
	for _, ref := range opts.Refs {
		if ref == "" {
			continue
		}
		args = append(args, ref)
	}
	return args
}
func filterRejected(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "filter") &&
		(strings.Contains(msg, "not allowed") ||
			strings.Contains(msg, "not recognized") ||
			strings.Contains(msg, "filtering") ||
			strings.Contains(msg, "invalid filter") ||
			strings.Contains(msg, "unknown option"))
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
