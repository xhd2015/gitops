package repos

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// EnsureStaticBareCommitsOptions configures EnsureStaticBareCommits.
type EnsureStaticBareCommitsOptions struct {
	// CloneURL is required: used as the git remote argument for FetchFromURL.
	CloneURL string
	Env      map[string]string

	// Primary is required when EnsureConnected is true (connectivity pivot).
	Primary string
	// EnsureConnected ensures Primary shares a merge-base with ConnectTo
	// commits (or every other commitIDs entry when ConnectTo is empty). If the
	// static bare is still shallow and not yet connected, EnsureStaticBareCommits
	// unshallows once under the repo flock.
	EnsureConnected bool
	// ConnectTo limits which commitIDs must share history with Primary when
	// EnsureConnected is true. Other commitIDs are still fetched for presence.
	// Empty means all commitIDs except Primary.
	ConnectTo []string
}

// EnsureStaticBareCommitsResult reports how the bare was filled.
type EnsureStaticBareCommitsResult struct {
	// Phase is fetch_commits | unshallow | connected.
	Phase string
	// LockWait is time spent blocked on the exclusive repo flock before critical section.
	LockWait time.Duration
}

// EnsureStaticBareCommits holds the exclusive repo flock once and:
//  1. fetches any missing commit SHAs (depth-1 FetchFromURL),
//  2. when EnsureConnected and histories are not yet linked, unshallows once
//     (blob:none) so the bare has full history — then connectivity is local.
//
// Unshallow-once replaces the old multi-step deepen ladder (16→32→…→unshallow).
// Peers wait on this flock for the holder's network work; a single unshallow
// keeps that critical section short, and later waiters usually see connected
// immediately. Tip blobs for worktrees still come from the depth-1 fetch in (1).
//
// Callers that only need missing SHAs may use FetchCommits; prefer this helper
// when EnsureConnected is needed so unshallow never races other fetches on the bare.
//
// FetchFromURL remains an unlocked primitive; this function is the safe multi-step
// entry for shared static-cache bares.
func EnsureStaticBareCommits(ctx context.Context, cacheDir string, commitIDs []string, opts EnsureStaticBareCommitsOptions) (EnsureStaticBareCommitsResult, error) {
	var result EnsureStaticBareCommitsResult
	result.Phase = "fetch_commits"

	if strings.TrimSpace(opts.CloneURL) == "" {
		return result, fmt.Errorf("EnsureStaticBareCommits: CloneURL is required")
	}
	primary := strings.TrimSpace(opts.Primary)
	if opts.EnsureConnected && primary == "" {
		return result, fmt.Errorf("EnsureStaticBareCommits: Primary is required when EnsureConnected")
	}

	lock, lockWait, err := acquireRepoLockContextTimed(ctx, cacheDir)
	result.LockWait = lockWait
	if err != nil {
		return result, err
	}
	defer releaseRepoLock(lock)

	if err := fetchCommitsLocked(ctx, cacheDir, commitIDs, FetchCommitsOptions{
		CloneURL: opts.CloneURL,
		Env:      opts.Env,
	}); err != nil {
		return result, err
	}

	if !opts.EnsureConnected {
		return result, nil
	}

	others := connectivityOthers(primary, commitIDs, opts.ConnectTo)
	if len(others) == 0 {
		result.Phase = "connected"
		return result, nil
	}

	ok, err := allHistoryConnected(ctx, cacheDir, primary, others)
	if err != nil {
		return result, err
	}
	if ok {
		result.Phase = "connected"
		return result, nil
	}

	if ctx.Err() != nil {
		result.Phase = "unshallow"
		return result, ctx.Err()
	}

	// One-shot full history under the flock (blob:none on FetchFromURL unshallow).
	// A peer that unshallowed while we waited is handled by the re-check below.
	result.Phase = "unshallow"
	_ = FetchFromURL(ctx, cacheDir, opts.CloneURL, URLFetchOptions{
		Unshallow: true,
		Env:       opts.Env,
	})
	ok, err = allHistoryConnected(ctx, cacheDir, primary, others)
	if err != nil {
		return result, err
	}
	if ok {
		return result, nil
	}
	return result, fmt.Errorf("EnsureStaticBareCommits: histories not connected after unshallow (primary=%s others=%v)", primary, others)
}

func allHistoryConnected(ctx context.Context, cacheDir, primary string, others []string) (bool, error) {
	for _, other := range others {
		connected, err := TargetedHistoryConnected(ctx, cacheDir, primary, other)
		if err != nil {
			return false, err
		}
		if !connected {
			return false, nil
		}
	}
	return true, nil
}

func connectivityOthers(primary string, commitIDs, connectTo []string) []string {
	src := commitIDs
	if len(connectTo) > 0 {
		src = connectTo
	}
	out := make([]string, 0, len(src))
	seen := map[string]bool{}
	for _, c := range src {
		c = strings.TrimSpace(c)
		if c == "" || c == primary || seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return out
}
