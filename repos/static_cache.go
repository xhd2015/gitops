package repos

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// EnsureStaticBareCache creates or reuses a shallow tip-only bare under
// {ReposRoot}/static-cache/<url-segments>/. Unlike EnsureBareCache with Depth>0
// (meta-only placeholder), this always materializes a real shallow bare so tip
// refs resolve. file:// remotes use --no-local so objects are not hardlinked
// from the full local history.
//
// Corrupt-bare recovery reuses the same holder-gated wipe logic as full cache:
// when BareHoldCount > 0, wipe is refused.
func EnsureStaticBareCache(ctx context.Context, cloneURL string, opts BareCacheOptions) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	cacheDir, err := StaticCacheDirUnder(opts.ReposRoot, cloneURL)
	if err != nil {
		return "", err
	}

	auth, gitURL, err := mergeAuth(opts.Auth, cloneURL)
	if err != nil {
		return "", fmt.Errorf("parse auth: %w", err)
	}
	gitURL = cloneURLForGit(gitURL, auth)

	if _, statErr := os.Stat(cacheDir); statErr == nil {
		if bareCacheFullyReady(cacheDir) {
			isBare, err := isBareRepository(cacheDir)
			if err != nil {
				return "", err
			}
			if !isBare {
				return "", fmt.Errorf("static cache path exists but is not a bare repository: %s", cacheDir)
			}
			return cacheDir, nil
		}
		return finishStaticBareCacheReady(ctx, cacheDir, gitURL, auth, opts)
	} else if !os.IsNotExist(statErr) {
		return "", statErr
	}

	return materializeStaticBareCache(ctx, cacheDir, gitURL, auth, opts)
}

func finishStaticBareCacheReady(ctx context.Context, cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) (string, error) {
	lock, err := acquireRepoLockContext(ctx, cacheDir)
	if err != nil {
		return "", err
	}
	defer releaseRepoLock(lock)

	if err := ensureStaticBareCachePreparedLocked(ctx, cacheDir, gitURL, auth, opts); err != nil {
		return "", err
	}
	return cacheDir, nil
}

// ensureStaticBareCachePreparedLocked mirrors ensureBareCachePreparedLocked's
// holder-gated wipe, but materializes a shallow tip bare (no full-history fetch).
func ensureStaticBareCachePreparedLocked(ctx context.Context, cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) error {
	if bareCacheFullyReady(cacheDir) {
		return nil
	}
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
				return fmt.Errorf("remove corrupt static bare cache: %w", rmErr)
			}
			return materializeStaticBareCacheLocked(ctx, cacheDir, gitURL, auth, opts)
		}
		if !isBare {
			return fmt.Errorf("static cache path exists but is not a bare repository: %s", cacheDir)
		}
		// Incomplete bare: finish shallow prepare and write marker (no full fetch).
		return prepareStaticBareAfterClone(ctx, cacheDir, auth, opts.Env)
	} else if !os.IsNotExist(statErr) {
		return statErr
	}

	return materializeStaticBareCacheLocked(ctx, cacheDir, gitURL, auth, opts)
}

func materializeStaticBareCache(ctx context.Context, cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) (string, error) {
	lock, err := acquireRepoLockContext(ctx, cacheDir)
	if err != nil {
		return "", err
	}
	defer releaseRepoLock(lock)

	if err := materializeStaticBareCacheLocked(ctx, cacheDir, gitURL, auth, opts); err != nil {
		return "", err
	}
	return cacheDir, nil
}

func materializeStaticBareCacheLocked(ctx context.Context, cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) error {
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
			return prepareStaticBareAfterClone(ctx, cacheDir, auth, opts.Env)
		}
		return fmt.Errorf("static cache path exists but is not a bare repository: %s", cacheDir)
	}

	if err := cloneStaticBareCache(ctx, cacheDir, gitURL, auth, opts); err != nil {
		return err
	}
	return nil
}

// cloneStaticBareCache performs a tip-only bare clone.
// --no-local forces file:// remotes to transfer objects instead of hardlinking
// the full local object store (critical for shallow size guarantees).
func cloneStaticBareCache(ctx context.Context, cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) error {
	depth := opts.Depth
	if depth <= 0 {
		depth = 1
	}
	base := []string{"clone", "--bare", "--no-local", fmt.Sprintf("--depth=%d", depth)}

	var lastErr error
	for _, branch := range []string{"main", "master"} {
		args := append(append([]string{}, base...), "--branch", branch, gitURL, cacheDir)
		if err := runGitWithEnvContext(ctx, "", auth, opts.Env, args...); err == nil {
			return prepareStaticBareAfterClone(ctx, cacheDir, auth, opts.Env)
		} else {
			lastErr = err
			if ctx.Err() != nil {
				return ctx.Err()
			}
			// Best-effort cleanup of partial clone dir before retry.
			_ = os.RemoveAll(cacheDir)
		}
	}
	args := append(append([]string{}, base...), gitURL, cacheDir)
	err := runGitWithEnvContext(ctx, "", auth, opts.Env, args...)
	if err == nil {
		return prepareStaticBareAfterClone(ctx, cacheDir, auth, opts.Env)
	}
	if lastErr != nil {
		return fmt.Errorf("clone static bare cache: %s", auth.maskError(lastErr))
	}
	return fmt.Errorf("clone static bare cache: %s", auth.maskError(err))
}

// prepareStaticBareAfterClone configures origin fetch refspec and ensures HEAD
// resolves without performing a full-history fetch (which would defeat shallow).
func prepareStaticBareAfterClone(ctx context.Context, cacheDir string, auth *GitAuthConfig, env map[string]string) error {
	_ = auth
	_ = env
	if err := setBareCacheOriginFetchConfig(cacheDir); err != nil {
		return err
	}
	// Prefer fixing HEAD from already-cloned tip refs; avoid full fetch.
	if err := ensureBareCacheHEAD(cacheDir); err != nil {
		// Fallback: shallow re-fetch default tips only (still depth-bounded).
		if fetchErr := shallowFetchOriginTips(ctx, cacheDir, auth, env); fetchErr != nil {
			return err
		}
		if err2 := ensureBareCacheHEAD(cacheDir); err2 != nil {
			return err2
		}
	}
	return writeBareCacheReadyMarker(cacheDir)
}

func shallowFetchOriginTips(ctx context.Context, cacheDir string, auth *GitAuthConfig, env map[string]string) error {
	for _, args := range [][]string{
		{"fetch", "--depth=1", "--prune", "origin", bareCacheOriginFetch},
		{"fetch", "--depth=1", "--prune", "origin"},
	} {
		if err := runGitWithEnvContext(ctx, cacheDir, auth, env, args...); err == nil {
			return nil
		} else if !strings.Contains(err.Error(), "couldn't find remote ref") {
			return fmt.Errorf("shallow fetch static bare origin tips: %w", err)
		}
	}
	return nil
}
