package repos

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/gitops/git"
)

// BareCacheOptions configures bare cache creation.
type BareCacheOptions struct {
	Depth     int // 0 = full clone, 1 = shallow clone
	Auth      *GitAuthConfig
	Env       map[string]string // optional proxy and other git env vars
	ReposRoot string            // empty → ReposRoot(); when set, CacheDir uses this root
}

// FetchOptions configures bare cache fetch behavior.
type FetchOptions struct {
	Auth      *GitAuthConfig
	Env       map[string]string
	ReposRoot string // empty → ReposRoot(); targeted cache uses this root
}

// bareCacheOriginFetch maps remote heads to refs/remotes/origin/* so fetch does
// not force-update refs/heads/* while another worktree has a branch checked out.
const bareCacheOriginFetch = "+refs/heads/*:refs/remotes/origin/*"

// bareCacheReadyMarkerName is written under the bare cache after successful post-clone prepare.
const bareCacheReadyMarkerName = ".repos-cache-ready"

func bareCacheReadyMarkerPath(cacheDir string) string {
	return filepath.Join(cacheDir, bareCacheReadyMarkerName)
}

func bareCacheReadyMarkerPresent(cacheDir string) bool {
	_, err := os.Stat(bareCacheReadyMarkerPath(cacheDir))
	return err == nil
}

// bareCacheFullyReady is true only when the ready marker exists and HEAD resolves
// to a commit. A marker alone is not enough (stale/broken HEAD must re-prepare).
func bareCacheFullyReady(cacheDir string) bool {
	return bareCacheReadyMarkerPresent(cacheDir) && bareCacheHEADResolves(cacheDir)
}

func writeBareCacheReadyMarker(cacheDir string) error {
	if err := os.WriteFile(bareCacheReadyMarkerPath(cacheDir), []byte{}, 0o644); err != nil {
		return fmt.Errorf("write bare cache ready marker: %w", err)
	}
	return nil
}

func clearBareCacheReadyMarker(cacheDir string) {
	_ = os.Remove(bareCacheReadyMarkerPath(cacheDir))
}

// EnsureBareCache creates or reuses a bare repository for the clone URL.
// A bare directory without .repos-cache-ready is not treated as ready: Ensure
// finishes prepare (under lock) before returning.
// When opts.ReposRoot is non-empty, the bare is placed under that root's cache/
// tree instead of the default ~/.repos/cache.
func EnsureBareCache(ctx context.Context, cloneURL string, opts BareCacheOptions) (string, error) {
	_ = ctx

	cacheDir, err := CacheDirUnder(opts.ReposRoot, cloneURL)
	if err != nil {
		return "", err
	}

	auth, gitURL, err := mergeAuth(opts.Auth, cloneURL)
	if err != nil {
		return "", fmt.Errorf("parse auth: %w", err)
	}
	gitURL = cloneURLForGit(gitURL, auth)

	if _, statErr := os.Stat(cacheDir); statErr == nil {
		// Fast path only when fully prepared (marker + resolvable HEAD). Mid-clone
		// concurrent callers and incomplete bares serialize under lock.
		if bareCacheFullyReady(cacheDir) {
			isBare, err := isBareRepository(cacheDir)
			if err != nil {
				return "", err
			}
			if !isBare {
				return "", fmt.Errorf("cache path exists but is not a bare repository: %s", cacheDir)
			}
			return cacheDir, nil
		}
		return finishBareCacheReady(cacheDir, gitURL, auth, opts)
	} else if !os.IsNotExist(statErr) {
		return "", statErr
	}

	if opts.Depth > 0 {
		if err := writeCloneMeta(cacheDir, cloneURL, opts); err != nil {
			return "", fmt.Errorf("write clone meta: %w", err)
		}
		return cacheDir, nil
	}

	return materializeBareCache(cacheDir, gitURL, auth, opts)
}

// finishBareCacheReady completes post-clone prepare for an existing bare that
// lacks the ready marker. Callers must not hold the repo lock.
func finishBareCacheReady(cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) (string, error) {
	lock, err := acquireRepoLock(cacheDir)
	if err != nil {
		return "", err
	}
	defer releaseRepoLock(lock)

	if err := ensureBareCachePreparedLocked(cacheDir, gitURL, auth, opts); err != nil {
		return "", err
	}
	return cacheDir, nil
}

// ensureBareCachePreparedLocked makes sure cacheDir is a fully prepared bare
// with ready marker and resolvable HEAD. Caller must hold the repo lock.
func ensureBareCachePreparedLocked(cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) error {
	if bareCacheFullyReady(cacheDir) {
		return nil
	}
	// Stale marker without usable HEAD: force re-prepare.
	if bareCacheReadyMarkerPresent(cacheDir) {
		clearBareCacheReadyMarker(cacheDir)
	}

	if _, statErr := os.Stat(cacheDir); statErr == nil {
		isBare, err := isBareRepository(cacheDir)
		if err != nil {
			// Corrupt / incomplete git dir: re-clone only when no live holders.
			if n := BareHoldCount(cacheDir); n > 0 {
				return fmt.Errorf("bare cache in use (%d holders); refuse wipe: %w", n, err)
			}
			if rmErr := os.RemoveAll(cacheDir); rmErr != nil {
				return fmt.Errorf("remove corrupt bare cache: %w", rmErr)
			}
			return materializeBareCacheLocked(cacheDir, gitURL, auth, opts)
		}
		if !isBare {
			return fmt.Errorf("cache path exists but is not a bare repository: %s", cacheDir)
		}
		// Incomplete bare: finish prepare and write marker.
		return prepareBareCacheAfterClone(cacheDir, auth, opts.Env)
	} else if !os.IsNotExist(statErr) {
		return statErr
	}

	return materializeBareCacheLocked(cacheDir, gitURL, auth, opts)
}

func materializeBareCache(cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) (string, error) {
	lock, err := acquireRepoLock(cacheDir)
	if err != nil {
		return "", err
	}
	defer releaseRepoLock(lock)

	if err := materializeBareCacheLocked(cacheDir, gitURL, auth, opts); err != nil {
		return "", err
	}
	return cacheDir, nil
}

func materializeBareCacheLocked(cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) error {
	if bareCacheFullyReady(cacheDir) {
		return nil
	}
	if bareCacheReadyMarkerPresent(cacheDir) {
		clearBareCacheReadyMarker(cacheDir)
	}

	if _, statErr := os.Stat(cacheDir); statErr == nil {
		isBare, err := isBareRepository(cacheDir)
		if err != nil {
			return err
		}
		if isBare {
			// Another caller may have cloned but not finished prepare yet (or
			// legacy bare without marker / unresolvable HEAD). Complete prepare.
			return prepareBareCacheAfterClone(cacheDir, auth, opts.Env)
		}
		return fmt.Errorf("cache path exists but is not a bare repository: %s", cacheDir)
	}

	if err := cloneBareCache(cacheDir, gitURL, auth, opts); err != nil {
		return err
	}
	removeCloneMeta(cacheDir)
	return nil
}

func cloneBareCache(cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) error {
	base := []string{"clone", "--bare"}
	if opts.Depth > 0 {
		base = append(base, fmt.Sprintf("--depth=%d", opts.Depth))
		var lastErr error
		for _, branch := range []string{"main", "master"} {
			args := append(append([]string{}, base...), "--branch", branch, gitURL, cacheDir)
			if err := runGitWithEnv("", auth, opts.Env, args...); err == nil {
				return prepareBareCacheAfterClone(cacheDir, auth, opts.Env)
			} else {
				lastErr = err
			}
		}
		args := append(append([]string{}, base...), gitURL, cacheDir)
		err := runGitWithEnv("", auth, opts.Env, args...)
		if err == nil {
			return prepareBareCacheAfterClone(cacheDir, auth, opts.Env)
		}
		if lastErr != nil {
			return fmt.Errorf("clone bare cache: %s", auth.maskError(lastErr))
		}
		return fmt.Errorf("clone bare cache: %s", auth.maskError(err))
	}

	args := append(append([]string{}, base...), gitURL, cacheDir)
	if err := runGitWithEnv("", auth, opts.Env, args...); err != nil {
		return fmt.Errorf("clone bare cache: %s", auth.maskError(err))
	}
	return prepareBareCacheAfterClone(cacheDir, auth, opts.Env)
}

func prepareBareCacheAfterClone(cacheDir string, auth *GitAuthConfig, env map[string]string) error {
	if err := setBareCacheOriginFetchConfig(cacheDir); err != nil {
		return err
	}
	if err := fetchBareCacheOriginRefs(cacheDir, auth, env); err != nil {
		return err
	}
	// Remotes created with git init --bare often leave HEAD at master while only
	// main was pushed; bare clones inherit that unresolvable HEAD. Fix before ready.
	if err := ensureBareCacheHEAD(cacheDir); err != nil {
		return err
	}
	// Ready marker only after prepare succeeds and HEAD resolves to a commit.
	return writeBareCacheReadyMarker(cacheDir)
}

func bareCacheHEADResolves(cacheDir string) bool {
	out, err := gitOutput(cacheDir, "rev-parse", "--verify", "HEAD^{commit}")
	return err == nil && strings.TrimSpace(out) != ""
}

// ensureBareCacheHEAD makes HEAD resolve to a commit in the bare cache.
// Prefer local heads, then origin remote-tracking refs (post-fetch).
func ensureBareCacheHEAD(cacheDir string) error {
	if bareCacheHEADResolves(cacheDir) {
		return nil
	}

	// Best-effort: let git discover origin's default branch into origin/HEAD.
	_ = runGit(cacheDir, nil, "remote", "set-head", "origin", "--auto")

	candidates := make([]string, 0, 16)
	// Prefer explicit common defaults first.
	candidates = append(candidates,
		"refs/heads/main",
		"refs/heads/master",
		"refs/remotes/origin/main",
		"refs/remotes/origin/master",
	)
	if originHead, err := gitOutput(cacheDir, "symbolic-ref", "-q", "refs/remotes/origin/HEAD"); err == nil {
		originHead = strings.TrimSpace(originHead)
		if originHead != "" {
			candidates = append([]string{originHead}, candidates...)
		}
	}
	if listed, err := gitOutput(cacheDir, "for-each-ref", "--format=%(refname)", "refs/heads/", "refs/remotes/origin/"); err == nil {
		for _, line := range strings.Split(listed, "\n") {
			ref := strings.TrimSpace(line)
			if ref == "" || ref == "refs/remotes/origin/HEAD" {
				continue
			}
			candidates = append(candidates, ref)
		}
	}

	seen := make(map[string]bool, len(candidates))
	for _, ref := range candidates {
		if seen[ref] {
			continue
		}
		seen[ref] = true
		commit, err := gitOutput(cacheDir, "rev-parse", "--verify", ref+"^{commit}")
		if err != nil || strings.TrimSpace(commit) == "" {
			continue
		}
		// Prefer a local branch for HEAD when only a remote-tracking tip exists.
		headRef := ref
		if strings.HasPrefix(ref, "refs/remotes/origin/") {
			branch := strings.TrimPrefix(ref, "refs/remotes/origin/")
			if branch != "" && branch != "HEAD" {
				local := "refs/heads/" + branch
				if _, localErr := gitOutput(cacheDir, "rev-parse", "--verify", local+"^{commit}"); localErr != nil {
					// Materialize local branch from remote-tracking tip so HEAD can
					// point at refs/heads/* (worktree-friendly bare layout).
					if updErr := runGit(cacheDir, nil, "update-ref", local, commit); updErr != nil {
						// Fall back to pointing HEAD at the remote-tracking ref.
						headRef = ref
					} else {
						headRef = local
					}
				} else {
					headRef = local
				}
			}
		}
		if err := runGit(cacheDir, nil, "symbolic-ref", "HEAD", headRef); err != nil {
			return fmt.Errorf("set bare cache HEAD to %s: %w", headRef, err)
		}
		if bareCacheHEADResolves(cacheDir) {
			return nil
		}
	}
	return fmt.Errorf("bare cache HEAD unresolvable after prepare in %s", cacheDir)
}

func fetchBareCacheOriginRefs(cacheDir string, auth *GitAuthConfig, env map[string]string) error {
	for _, args := range [][]string{
		{"fetch", "--prune", "origin", bareCacheOriginFetch},
		{"fetch", "--prune", "origin"},
	} {
		if err := runGitWithEnv(cacheDir, auth, env, args...); err == nil {
			return nil
		} else if !strings.Contains(err.Error(), "couldn't find remote ref") {
			return fmt.Errorf("fetch bare cache origin refs: %w", err)
		}
	}
	return nil
}

func setBareCacheOriginFetchConfig(cacheDir string) error {
	if err := runGit(cacheDir, nil, "config", "remote.origin.fetch", bareCacheOriginFetch); err != nil {
		return fmt.Errorf("set bare cache fetch config: %w", err)
	}
	return nil
}

func bareCacheFetchRefspecNeedsMigration(refspec string) bool {
	if refspec == "" {
		return true
	}
	parts := strings.Split(refspec, ":")
	if len(parts) < 2 {
		return false
	}
	dest := parts[len(parts)-1]
	return strings.HasPrefix(dest, "refs/heads/")
}

func ensureBareCacheFetchConfig(cacheDir string) error {
	current, err := gitOutput(cacheDir, "config", "--get", "remote.origin.fetch")
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		return fmt.Errorf("read bare cache fetch config: %w", err)
	}
	if !bareCacheFetchRefspecNeedsMigration(current) {
		return nil
	}
	return setBareCacheOriginFetchConfig(cacheDir)
}

func ensureBareCacheReadyLocked(cacheDir string, auth *GitAuthConfig) error {
	isBare, err := isBareRepository(cacheDir)
	if err == nil && isBare {
		if bareCacheFullyReady(cacheDir) {
			return nil
		}
		if bareCacheReadyMarkerPresent(cacheDir) {
			clearBareCacheReadyMarker(cacheDir)
		}
		// Legacy/incomplete bare without ready marker or with broken HEAD.
		return prepareBareCacheAfterClone(cacheDir, auth, nil)
	}
	if err == nil {
		return fmt.Errorf("cache path exists but is not a bare repository: %s", cacheDir)
	}

	meta, err := readCloneMeta(cacheDir)
	if err != nil {
		return fmt.Errorf("bare cache not initialized for %s: %w", cacheDir, err)
	}
	optAuth := auth
	if optAuth == nil {
		optAuth, _, err = mergeAuth(nil, meta.CloneURL)
		if err != nil {
			return err
		}
	}
	gitURL := cloneURLForGit(meta.CloneURL, optAuth)
	return materializeBareCacheLocked(cacheDir, gitURL, optAuth, BareCacheOptions{Depth: meta.Depth, Auth: optAuth})
}

// FetchBareCache fetches updates for an existing bare cache under file lock.
func FetchBareCache(ctx context.Context, cacheDir string, opts FetchOptions) error {
	lock, err := acquireRepoLockContext(ctx, cacheDir)
	if err != nil {
		return err
	}
	defer releaseRepoLock(lock)

	if err := ensureBareCacheReadyLocked(cacheDir, opts.Auth); err != nil {
		return err
	}

	if err := ensureBareCacheFetchConfig(cacheDir); err != nil {
		return err
	}

	return fetchBareCacheOriginRefsContext(ctx, cacheDir, opts.Auth, opts.Env)
}

// FetchBareCacheRefs refreshes only the requested remote branch refs. It is
// intended for latency-sensitive callers that do not need a complete mirror.
func FetchBareCacheRefs(ctx context.Context, cacheDir string, refs []string, opts FetchOptions) error {
	branches := fetchableBranchRefs(refs)
	if len(branches) == 0 {
		return nil
	}
	lock, err := acquireRepoLockContext(ctx, cacheDir)
	if err != nil {
		return err
	}
	defer releaseRepoLock(lock)
	if err := ensureBareCacheReadyLocked(cacheDir, opts.Auth); err != nil {
		return err
	}
	args := []string{"fetch", "--no-tags", "origin"}
	for _, branch := range branches {
		args = append(args, "+refs/heads/"+branch+":refs/remotes/origin/"+branch)
	}
	if err := runGitWithEnvContext(ctx, cacheDir, opts.Auth, opts.Env, args...); err != nil {
		return fmt.Errorf("fetch requested branch refs: %w", err)
	}
	return nil
}

func fetchableBranchRefs(refs []string) []string {
	seen := make(map[string]bool, len(refs))
	branches := make([]string, 0, len(refs))
	for _, ref := range refs {
		branch := ""
		switch {
		case strings.HasPrefix(ref, "refs/heads/"):
			branch = strings.TrimPrefix(ref, "refs/heads/")
		case strings.HasPrefix(ref, "refs/remotes/origin/"):
			branch = strings.TrimPrefix(ref, "refs/remotes/origin/")
		case strings.HasPrefix(ref, "origin/"):
			branch = strings.TrimPrefix(ref, "origin/")
		}
		if branch == "" || seen[branch] {
			continue
		}
		seen[branch] = true
		branches = append(branches, branch)
	}
	sort.Strings(branches)
	return branches
}

// FetchCommitsOptions supplies credentials/proxy for deepen.
// When CloneURL is set, FetchCommits uses FetchFromURL (URL-as-remote) instead of bare "origin".
type FetchCommitsOptions struct {
	CloneURL string
	Env      map[string]string
}

// FetchCommits deepens a shallow bare cache to include the given commits.
// Existing callers keep compiling as FetchCommits(ctx, dir, ids).
// When opts provide CloneURL, uses FetchFromURL (depth 1) instead of bare "origin".
//
// Holds the exclusive repo flock for the whole missing-commit fetch so concurrent
// static/partial acquires on the same bare cannot race on shallow.lock.
// For fetch + connectivity deepen in one critical section, use EnsureStaticBareCommits.
func FetchCommits(ctx context.Context, cacheDir string, commitIDs []string, opts ...FetchCommitsOptions) error {
	var o FetchCommitsOptions
	if len(opts) > 0 {
		o = opts[0]
	}

	lock, err := acquireRepoLockContext(ctx, cacheDir)
	if err != nil {
		return err
	}
	defer releaseRepoLock(lock)

	return fetchCommitsLocked(ctx, cacheDir, commitIDs, o)
}

// fetchCommitsLocked is the body of FetchCommits; caller must hold the repo flock.
// Non-reentrant: do not call FetchCommits from under the same flock.
func fetchCommitsLocked(ctx context.Context, cacheDir string, commitIDs []string, o FetchCommitsOptions) error {
	if err := ensureBareCacheReadyLocked(cacheDir, nil); err != nil {
		return err
	}

	missing := make([]string, 0, len(commitIDs))
	for _, commitID := range commitIDs {
		if commitID == "" {
			continue
		}
		if resolved, _ := git.RevParseOrEmpty(cacheDir, commitID); resolved != "" {
			continue
		}
		missing = append(missing, commitID)
	}
	if len(missing) == 0 {
		return nil
	}

	if o.CloneURL != "" {
		// Unified path: clone URL as remote argument + Env (proxy); no extraHeader.
		if err := FetchFromURL(ctx, cacheDir, o.CloneURL, URLFetchOptions{
			Refs:  missing,
			Depth: 1,
			Env:   o.Env,
		}); err != nil {
			return err
		}
		return nil
	}

	// Legacy/sibling path: deepen via configured origin remote.
	for _, commitID := range missing {
		if err := runGitWithEnvContext(ctx, cacheDir, nil, o.Env, "fetch", "--depth=1", "origin", commitID); err != nil {
			return fmt.Errorf("fetch commit %s: %w", commitID, err)
		}
	}
	return nil
}

func fetchBareCacheOriginRefsContext(ctx context.Context, cacheDir string, auth *GitAuthConfig, env map[string]string) error {
	for _, args := range [][]string{
		{"fetch", "--prune", "origin", bareCacheOriginFetch},
		{"fetch", "--prune", "origin"},
	} {
		if err := runGitWithEnvContext(ctx, cacheDir, auth, env, args...); err == nil {
			return nil
		} else if !strings.Contains(err.Error(), "couldn't find remote ref") {
			return fmt.Errorf("fetch bare cache origin refs: %w", err)
		}
	}
	return nil
}
