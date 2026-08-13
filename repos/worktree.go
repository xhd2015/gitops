package repos

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WorktreeEntry describes a git worktree registered with a bare cache.
type WorktreeEntry struct {
	Path   string
	HEAD   string
	Branch string // e.g. refs/heads/master; empty when detached
}

// ListWorktrees returns worktrees registered for a bare cache repository.
func ListWorktrees(cacheDir string) ([]WorktreeEntry, error) {
	out, err := gitOutput(cacheDir, "worktree", "list", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("list worktrees: %w", err)
	}
	if out == "" {
		return nil, nil
	}

	var entries []WorktreeEntry
	var current WorktreeEntry
	flush := func() {
		if current.Path == "" {
			return
		}
		entries = append(entries, current)
		current = WorktreeEntry{}
	}

	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			flush()
			continue
		}
		fields := strings.SplitN(line, " ", 2)
		if len(fields) != 2 {
			continue
		}
		switch fields[0] {
		case "worktree":
			flush()
			current.Path = fields[1]
		case "HEAD":
			current.HEAD = fields[1]
		case "branch":
			current.Branch = fields[1]
		}
	}
	flush()
	return entries, nil
}

// PruneConflictingWorktrees removes stale worktrees that block fetch or re-checkout.
// It keeps a worktree at keepPath when it already tracks the same branch.
func PruneConflictingWorktrees(cacheDir, ref string, isBranch bool, keepPath string) error {
	entries, err := ListWorktrees(cacheDir)
	if err != nil {
		return err
	}

	branchRef := ""
	if isBranch {
		branchRef = "refs/heads/" + ref
	}

	cacheDir = filepath.Clean(cacheDir)
	keepPath = filepath.Clean(keepPath)
	for _, entry := range entries {
		entryPath := filepath.Clean(entry.Path)
		if entryPath == cacheDir {
			continue
		}

		remove := false
		if strings.HasPrefix(entryPath, cacheDir+string(os.PathSeparator)) {
			remove = true
		}
		if isBranch && entry.Branch == branchRef && entryPath != keepPath {
			remove = true
		}
		if !remove {
			continue
		}
		if err := RemoveWorktree(cacheDir, entryPath); err != nil {
			return err
		}
	}
	return nil
}

// FetchRemoteBranch updates refs/remotes/origin/<branch> without touching refs/heads.
func FetchRemoteBranch(cacheDir, branch string) error {
	refspec := fmt.Sprintf("refs/heads/%s:refs/remotes/origin/%s", branch, branch)
	if err := runGit(cacheDir, nil, "fetch", "origin", refspec); err != nil {
		return fmt.Errorf("fetch remote branch %s: %w", branch, err)
	}
	return nil
}

// AddWorktree creates a worktree from a bare cache at the given ref or commit.
// It is idempotent: an existing complete worktree (.git present) is a no-op.
// Incomplete directories are cleaned and recreated. Concurrent adds that race
// on "already exists" / "already checked out" / locked-worktree succeed when
// .git is present, otherwise unlock/prune/remove incomplete and retry once
// with force (branch-tracking, not detach).
func AddWorktree(cacheDir, worktreeDir, ref string) error {
	if err := os.MkdirAll(filepath.Dir(worktreeDir), 0o755); err != nil {
		return fmt.Errorf("create worktree parent dir: %w", err)
	}

	gitMarker := filepath.Join(worktreeDir, ".git")
	// Complete worktree already present — idempotent success.
	if _, err := os.Stat(gitMarker); err == nil {
		return nil
	}

	// Incomplete path (exists but no .git): drop stale registration and dir.
	if _, err := os.Stat(worktreeDir); err == nil {
		clearStaleWorktree(cacheDir, worktreeDir)
	}

	// Best-effort: free the branch if checked out elsewhere under this cache.
	_ = PruneConflictingWorktrees(cacheDir, ref, true, worktreeDir)

	if err := runGit(cacheDir, nil, "worktree", "add", worktreeDir, ref); err != nil {
		if !isRecoverableWorktreeAddErr(err) {
			return fmt.Errorf("add worktree: %w", err)
		}
		// Concurrent winner may have finished while we failed.
		if _, statErr := os.Stat(gitMarker); statErr == nil {
			return nil
		}
		clearStaleWorktree(cacheDir, worktreeDir)
		if _, statErr := os.Stat(gitMarker); statErr == nil {
			return nil
		}
		_ = PruneConflictingWorktrees(cacheDir, ref, true, worktreeDir)
		// Force retry clears "missing but locked" admin records (git: add -f -f).
		if retryErr := runGit(cacheDir, nil, "worktree", "add", "-f", "-f", worktreeDir, ref); retryErr != nil {
			if _, statErr := os.Stat(gitMarker); statErr == nil {
				return nil
			}
			// Second recovery pass if race re-locked between clear and retry.
			if isRecoverableWorktreeAddErr(retryErr) {
				clearStaleWorktree(cacheDir, worktreeDir)
				if _, statErr := os.Stat(gitMarker); statErr == nil {
					return nil
				}
				if retry2 := runGit(cacheDir, nil, "worktree", "add", "-f", "-f", worktreeDir, ref); retry2 != nil {
					if _, statErr := os.Stat(gitMarker); statErr == nil {
						return nil
					}
					return fmt.Errorf("add worktree: %w", retry2)
				}
				return nil
			}
			return fmt.Errorf("add worktree: %w", retryErr)
		}
		return nil
	}
	return nil
}

func isRecoverableWorktreeAddErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "already exists") ||
		strings.Contains(msg, "already checked out") ||
		strings.Contains(msg, "locked worktree") ||
		strings.Contains(msg, "missing but locked")
}

// clearStaleWorktree unlocks, prunes, and removes a partial/concurrent worktree registration.
// If a concurrent peer finished (.git present) after unlock/prune, leaves it alone.
func clearStaleWorktree(cacheDir, worktreeDir string) {
	_ = runGit(cacheDir, nil, "worktree", "unlock", worktreeDir)
	_ = runGit(cacheDir, nil, "worktree", "prune")
	// Peer may have completed while we unlocked/pruned — do not delete a valid checkout.
	if _, err := os.Stat(filepath.Join(worktreeDir, ".git")); err == nil {
		return
	}
	_ = runGit(cacheDir, nil, "worktree", "remove", "-f", worktreeDir)
	_ = os.RemoveAll(worktreeDir)
}

// AddDetachedWorktree creates a detached worktree from a bare cache.
// It is safe for private (unique-per-caller) paths: a complete worktree
// (.git present) is left alone (idempotent); incomplete directories are
// cleared, then `worktree add -f --detach` runs. On recoverable "already
// exists" races for that private path, succeeds if .git is present, else
// clears and retries once. Does not implement refcount shared-path reuse.
func AddDetachedWorktree(cacheDir, worktreeDir, ref string) error {
	if err := os.MkdirAll(filepath.Dir(worktreeDir), 0o755); err != nil {
		return fmt.Errorf("create worktree parent dir: %w", err)
	}

	gitMarker := filepath.Join(worktreeDir, ".git")
	// Complete worktree already present — idempotent success for this private path.
	if _, err := os.Stat(gitMarker); err == nil {
		return nil
	}

	// Incomplete path (exists but no .git): drop stale registration and dir.
	if _, err := os.Stat(worktreeDir); err == nil {
		clearStaleWorktree(cacheDir, worktreeDir)
	}

	if err := runGit(cacheDir, nil, "worktree", "add", "-f", "--detach", worktreeDir, ref); err != nil {
		if !isRecoverableWorktreeAddErr(err) {
			return fmt.Errorf("add detached worktree: %w", err)
		}
		// Concurrent winner may have finished while we failed.
		if _, statErr := os.Stat(gitMarker); statErr == nil {
			return nil
		}
		clearStaleWorktree(cacheDir, worktreeDir)
		if _, statErr := os.Stat(gitMarker); statErr == nil {
			return nil
		}
		if retryErr := runGit(cacheDir, nil, "worktree", "add", "-f", "--detach", worktreeDir, ref); retryErr != nil {
			if _, statErr := os.Stat(gitMarker); statErr == nil {
				return nil
			}
			return fmt.Errorf("add detached worktree: %w", retryErr)
		}
		return nil
	}
	return nil
}

// AddSparseDetachedWorktree creates a detached worktree with non-cone
// sparse-checkout of paths (empty paths are allowed). Does not use --cone.
func AddSparseDetachedWorktree(cacheDir, worktreeDir, sha string, paths []string) error {
	if err := os.MkdirAll(filepath.Dir(worktreeDir), 0o755); err != nil {
		return fmt.Errorf("create worktree parent dir: %w", err)
	}

	if _, err := os.Stat(worktreeDir); err == nil {
		clearStaleWorktree(cacheDir, worktreeDir)
	}

	if err := runGit(cacheDir, nil, "worktree", "add", "--no-checkout", "-f", "--detach", worktreeDir, sha); err != nil {
		return fmt.Errorf("add sparse detached worktree: %w", err)
	}
	if err := runGit(worktreeDir, nil, "sparse-checkout", "init", "--no-cone"); err != nil {
		_ = RemoveWorktree(cacheDir, worktreeDir)
		return fmt.Errorf("sparse-checkout init: %w", err)
	}
	setArgs := append([]string{"sparse-checkout", "set", "--no-cone", "--"}, paths...)
	if err := runGit(worktreeDir, nil, setArgs...); err != nil {
		_ = RemoveWorktree(cacheDir, worktreeDir)
		return fmt.Errorf("sparse-checkout set: %w", err)
	}
	if err := runGit(worktreeDir, nil, "checkout", "--force", sha); err != nil {
		if err2 := runGit(worktreeDir, nil, "read-tree", "-mu", "HEAD"); err2 != nil {
			_ = RemoveWorktree(cacheDir, worktreeDir)
			return fmt.Errorf("sparse checkout materialize: %w", err)
		}
	}
	return nil
}

// RemoveWorktree removes a worktree from a bare cache.
func RemoveWorktree(cacheDir, worktreeDir string) error {
	if _, err := os.Stat(worktreeDir); os.IsNotExist(err) {
		return nil
	}
	if err := runGit(cacheDir, nil, "worktree", "remove", "-f", worktreeDir); err != nil {
		return fmt.Errorf("remove worktree: %w", err)
	}
	return nil
}

// AddPersistentWorktree ensures bare cache and creates a persistent worktree.
func AddPersistentWorktree(ctx context.Context, cloneURL, ref string) (string, error) {
	cacheDir, err := EnsureBareCache(ctx, cloneURL, BareCacheOptions{Depth: 0})
	if err != nil {
		return "", err
	}
	worktreeDir, err := WorktreeDir(cloneURL, ref)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(worktreeDir); err == nil {
		return worktreeDir, nil
	}
	if err := AddWorktree(cacheDir, worktreeDir, ref); err != nil {
		return "", err
	}
	return worktreeDir, nil
}

// AddTmpWorktree ensures bare cache and creates an ephemeral worktree.
func AddTmpWorktree(ctx context.Context, cloneURL, id, ref string) (string, error) {
	cacheDir, err := EnsureBareCache(ctx, cloneURL, BareCacheOptions{Depth: 0})
	if err != nil {
		return "", err
	}
	worktreeDir, err := TmpWorktreeDir(cloneURL, id)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(worktreeDir); err == nil {
		return worktreeDir, nil
	}
	if err := AddWorktree(cacheDir, worktreeDir, ref); err != nil {
		return "", err
	}
	return worktreeDir, nil
}
