package repos

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhd2015/gitops/git"
)

// EnsureCommitsOptions configures commit-only bare cache materialization.
type EnsureCommitsOptions struct {
	Auth *GitAuthConfig
	Env  map[string]string
	// ReposRoot injects the unified repos root (empty → ReposRoot() / real home).
	ReposRoot string
	// Progress, when set, prints dir + exact git argv and streams fetch --progress.
	Progress io.Writer
}

// EnsureCommits ensures the bare cache for cloneURL can resolve every commitID.
// On a cold cache it bootstraps with git init --bare + commit-only fetches
// (git fetch --depth=1 origin <sha>), without probing main/master or bulk-fetching
// remote heads. On a warm cache it fetches only missing SHAs.
func EnsureCommits(ctx context.Context, cloneURL string, commitIDs []string, opts EnsureCommitsOptions) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if len(commitIDs) == 0 {
		return "", fmt.Errorf("commitIDs must not be empty")
	}

	cacheDir, err := CacheDirUnder(opts.ReposRoot, cloneURL)
	if err != nil {
		return "", err
	}

	auth, gitURL, err := mergeAuth(opts.Auth, cloneURL)
	if err != nil {
		return "", fmt.Errorf("parse auth: %w", err)
	}
	gitURL = cloneURLForGit(gitURL, auth)

	lock, err := acquireRepoLockContext(ctx, cacheDir)
	if err != nil {
		return "", err
	}
	defer releaseRepoLock(lock)

	if err := ensureBareOriginForCommits(cacheDir, gitURL, auth, opts.Env); err != nil {
		return "", err
	}

	for _, commitID := range commitIDs {
		if commitID == "" {
			continue
		}
		if resolved, _ := git.RevParseOrEmpty(cacheDir, commitID); resolved != "" {
			if opts.Progress != nil {
				fmt.Fprintf(opts.Progress, "have %s\n", commitID)
			}
			continue
		}
		fetchArgs := []string{"fetch", "--depth=1", "origin", commitID}
		if opts.Progress != nil {
			fetchArgs = []string{"fetch", "--progress", "--depth=1", "origin", commitID}
			fmt.Fprintf(opts.Progress, "dir: %s\n", cacheDir)
			fmt.Fprintf(opts.Progress, "$ git %s\n", strings.Join(fetchArgs, " "))
			if err := runGitStreamStderr(ctx, cacheDir, auth, opts.Env, opts.Progress, fetchArgs...); err != nil {
				return "", fmt.Errorf("fetch commit %s: %s", commitID, auth.maskError(err))
			}
		} else if err := runGitWithEnvContext(ctx, cacheDir, auth, opts.Env, fetchArgs...); err != nil {
			return "", fmt.Errorf("fetch commit %s: %s", commitID, auth.maskError(err))
		}
		if resolved, _ := git.RevParseOrEmpty(cacheDir, commitID); resolved == "" {
			return "", fmt.Errorf("fetch commit %s: object still not resolvable after fetch", commitID)
		}
	}

	return cacheDir, nil
}

// ensureBareOriginForCommits makes sure cacheDir is a bare repo with origin set
// to gitURL. It does not fetch remote heads.
func ensureBareOriginForCommits(cacheDir, gitURL string, auth *GitAuthConfig, env map[string]string) error {
	if _, statErr := os.Stat(cacheDir); statErr == nil {
		isBare, err := isBareRepository(cacheDir)
		if err != nil {
			return err
		}
		if !isBare {
			return fmt.Errorf("cache path exists but is not a bare repository: %s", cacheDir)
		}
		return ensureOriginRemote(cacheDir, gitURL, auth, env)
	} else if !os.IsNotExist(statErr) {
		return statErr
	}

	if err := os.MkdirAll(filepath.Dir(cacheDir), 0o755); err != nil {
		return fmt.Errorf("create cache parent dir: %w", err)
	}
	if err := runGitWithEnv("", auth, env, "init", "--bare", cacheDir); err != nil {
		return fmt.Errorf("init bare cache: %s", auth.maskError(err))
	}
	if err := runGitWithEnv(cacheDir, auth, env, "remote", "add", "origin", gitURL); err != nil {
		return fmt.Errorf("add origin remote: %s", auth.maskError(err))
	}
	// Intentionally do not set +refs/heads/*:refs/remotes/origin/* or bulk-fetch.
	removeCloneMeta(cacheDir)
	return nil
}

func ensureOriginRemote(cacheDir, gitURL string, auth *GitAuthConfig, env map[string]string) error {
	out, err := gitOutput(cacheDir, "remote", "get-url", "origin")
	if err == nil {
		// Origin exists; update URL if needed so auth/path changes still work.
		if strings.TrimSpace(out) != gitURL && gitURL != "" {
			if setErr := runGitWithEnv(cacheDir, auth, env, "remote", "set-url", "origin", gitURL); setErr != nil {
				return fmt.Errorf("set origin url: %s", auth.maskError(setErr))
			}
		}
		return nil
	}
	if err := runGitWithEnv(cacheDir, auth, env, "remote", "add", "origin", gitURL); err != nil {
		return fmt.Errorf("add origin remote: %s", auth.maskError(err))
	}
	return nil
}
