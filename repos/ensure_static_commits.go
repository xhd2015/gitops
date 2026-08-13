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
	// EnsureConnected grows shallow history until Primary shares a merge-base
	// with every other commit (or unshallow is reached).
	EnsureConnected bool
}

// EnsureStaticBareCommitsResult reports how the bare was filled.
type EnsureStaticBareCommitsResult struct {
	// Phase is fetch_commits | deepen_16 | deepen_32 | deepen_64 | deepen_256 |
	// deepen_1024 | unshallow | connected.
	Phase string
	// LockWait is time spent blocked on the exclusive repo flock before critical section.
	LockWait time.Duration
}

type staticDeepenPhase struct {
	name      string
	deepen    int
	unshallow bool
}

// Same ladder as targeted-ref connectivity (must stay in lock for the whole climb).
// git --deepen=N is incremental from the current shallow boundary; start small so
// near-tip merges stop early without paying a full +64 fetch first.
var staticConnectivityPhases = []staticDeepenPhase{
	{name: "deepen_16", deepen: 16},
	{name: "deepen_32", deepen: 32},
	{name: "deepen_64", deepen: 64},
	{name: "deepen_256", deepen: 256},
	{name: "deepen_1024", deepen: 1024},
	{name: "unshallow", unshallow: true},
}

// EnsureStaticBareCommits holds the exclusive repo flock once and:
//  1. fetches any missing commit SHAs (depth-1 FetchFromURL),
//  2. optionally deepens until Primary connects with every other commit,
//     rechecking connectivity after lock wait and before each deepen step
//     so concurrent waiters skip redundant network work.
//
// Callers that only need missing SHAs may use FetchCommits; prefer this helper
// when EnsureConnected is needed so deepen never races other fetches on the bare.
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

	others := make([]string, 0, len(commitIDs))
	for _, c := range commitIDs {
		c = strings.TrimSpace(c)
		if c == "" || c == primary {
			continue
		}
		others = append(others, c)
	}
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

	// Progressive history growth under the same flock (no unlock between steps).
	for _, p := range staticConnectivityPhases {
		if ctx.Err() != nil {
			result.Phase = p.name
			return result, ctx.Err()
		}
		// Recheck: another waiter may have finished this phase while we queued.
		ok, err = allHistoryConnected(ctx, cacheDir, primary, others)
		if err != nil {
			result.Phase = p.name
			return result, err
		}
		if ok {
			result.Phase = "connected"
			return result, nil
		}

		fetchOpts := URLFetchOptions{
			Refs:      commitIDs,
			Deepen:    p.deepen,
			Unshallow: p.unshallow,
			Env:       opts.Env,
		}
		if ferr := FetchFromURL(ctx, cacheDir, opts.CloneURL, fetchOpts); ferr != nil {
			// Unshallow on a non-shallow repo or deepen no-op: still re-check connect.
			if !p.unshallow {
				result.Phase = p.name
				return result, fmt.Errorf("EnsureStaticBareCommits %s: %w", p.name, ferr)
			}
		}
		ok, err = allHistoryConnected(ctx, cacheDir, primary, others)
		if err != nil {
			result.Phase = p.name
			return result, err
		}
		if ok {
			result.Phase = p.name
			return result, nil
		}
		if p.unshallow {
			result.Phase = p.name
			return result, fmt.Errorf("EnsureStaticBareCommits: histories not connected after unshallow (primary=%s others=%v)", primary, others)
		}
	}
	return result, fmt.Errorf("EnsureStaticBareCommits: histories not connected (primary=%s)", primary)
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
