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
	_ = ctx

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
		return finishStaticBareCacheReady(cacheDir, gitURL, auth, opts)
	} else if !os.IsNotExist(statErr) {
		return "", statErr
	}

	return materializeStaticBareCache(cacheDir, gitURL, auth, opts)
}

func finishStaticBareCacheReady(cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) (string, error) {
	lock, err := acquireRepoLock(cacheDir)
	if err != nil {
		return "", err
	}
	defer releaseRepoLock(lock)

	if err := ensureStaticBareCachePreparedLocked(cacheDir, gitURL, auth, opts); err != nil {
		return "", err
	}
	return cacheDir, nil
}

// ensureStaticBareCachePreparedLocked mirrors ensureBareCachePreparedLocked's
// holder-gated wipe, but materializes a shallow tip bare (no full-history fetch).
func ensureStaticBareCachePreparedLocked(cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) error {
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
			return materializeStaticBareCacheLocked(cacheDir, gitURL, auth, opts)
		}
		if !isBare {
			return fmt.Errorf("static cache path exists but is not a bare repository: %s", cacheDir)
		}
		// Incomplete bare: finish shallow prepare and write marker (no full fetch).
		return prepareStaticBareAfterClone(cacheDir, auth, opts.Env)
	} else if !os.IsNotExist(statErr) {
		return statErr
	}

	return materializeStaticBareCacheLocked(cacheDir, gitURL, auth, opts)
}

func materializeStaticBareCache(cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) (string, error) {
	lock, err := acquireRepoLock(cacheDir)
	if err != nil {
		return "", err
	}
	defer releaseRepoLock(lock)

	if err := materializeStaticBareCacheLocked(cacheDir, gitURL, auth, opts); err != nil {
		return "", err
	}
	return cacheDir, nil
}

func materializeStaticBareCacheLocked(cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) error {
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
			return prepareStaticBareAfterClone(cacheDir, auth, opts.Env)
		}
		return fmt.Errorf("static cache path exists but is not a bare repository: %s", cacheDir)
	}

	if err := cloneStaticBareCache(cacheDir, gitURL, auth, opts); err != nil {
		return err
	}
	return nil
}

// cloneStaticBareCache performs a tip-only bare clone.
// --no-local forces file:// remotes to transfer objects instead of hardlinking
// the full local object store (critical for shallow size guarantees).
func cloneStaticBareCache(cacheDir, gitURL string, auth *GitAuthConfig, opts BareCacheOptions) error {
	depth := opts.Depth
	if depth <= 0 {
		depth = 1
	}
	base := []string{"clone", "--bare", "--no-local", fmt.Sprintf("--depth=%d", depth)}

	var lastErr error
	for _, branch := range []string{"main", "master"} {
		args := append(append([]string{}, base...), "--branch", branch, gitURL, cacheDir)
		if err := runGitWithEnv("", auth, opts.Env, args...); err == nil {
			return prepareStaticBareAfterClone(cacheDir, auth, opts.Env)
		} else {
			lastErr = err
			// Best-effort cleanup of partial clone dir before retry.
			_ = os.RemoveAll(cacheDir)
		}
	}
	args := append(append([]string{}, base...), gitURL, cacheDir)
	err := runGitWithEnv("", auth, opts.Env, args...)
	if err == nil {
		return prepareStaticBareAfterClone(cacheDir, auth, opts.Env)
	}
	if lastErr != nil {
		return fmt.Errorf("clone static bare cache: %s", auth.maskError(lastErr))
	}
	return fmt.Errorf("clone static bare cache: %s", auth.maskError(err))
}

// prepareStaticBareAfterClone configures origin fetch refspec and ensures HEAD
// resolves without performing a full-history fetch (which would defeat shallow).
func prepareStaticBareAfterClone(cacheDir string, auth *GitAuthConfig, env map[string]string) error {
	_ = auth
	_ = env
	if err := setBareCacheOriginFetchConfig(cacheDir); err != nil {
		return err
	}
	// Prefer fixing HEAD from already-cloned tip refs; avoid full fetch.
	if err := ensureBareCacheHEAD(cacheDir); err != nil {
		// Fallback: shallow re-fetch default tips only (still depth-bounded).
		if fetchErr := shallowFetchOriginTips(cacheDir, auth, env); fetchErr != nil {
			return err
		}
		if err2 := ensureBareCacheHEAD(cacheDir); err2 != nil {
			return err2
		}
	}
	return writeBareCacheReadyMarker(cacheDir)
}

func shallowFetchOriginTips(cacheDir string, auth *GitAuthConfig, env map[string]string) error {
	for _, args := range [][]string{
		{"fetch", "--depth=1", "--prune", "origin", bareCacheOriginFetch},
		{"fetch", "--depth=1", "--prune", "origin"},
	} {
		if err := runGitWithEnv(cacheDir, auth, env, args...); err == nil {
			return nil
		} else if !strings.Contains(err.Error(), "couldn't find remote ref") {
			return fmt.Errorf("shallow fetch static bare origin tips: %w", err)
		}
	}
	return nil
}
