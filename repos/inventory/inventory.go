// Package inventory provides a pure, injectable git-cache inventory collector.
//
// Collect walks unified layout directories under an injectable ReposRoot (and
// optionally a tmp-worktrees root) and returns a Snapshot of sizes and counts.
// There is no Prometheus / metrics SDK dependency and no network I/O.
package inventory

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/xhd2015/gitops/repos"
)

// Layout path-segment keys (must match repos path helpers).
const (
	LayoutStaticCache   = "static-cache"
	LayoutWorktrees     = "worktrees"
	LayoutCache         = "cache"
	LayoutTargetedCache = "targeted-cache"
)

// Options configures Collect. Empty ReposRoot → repos.ReposRoot(); tests inject.
type Options struct {
	ReposRoot           string
	IncludeTmpWorktrees bool
	TmpWorktreesRoot    string        // empty + IncludeTmpWorktrees → repos.TmpWorktreesRoot()
	MaxWalkDuration     time.Duration // 0 = no soft limit
}

// LayoutStats is du-style stats for one layout directory (or tmp root).
type LayoutStats struct {
	Present       bool
	Bytes         int64
	RepoCount     int64
	WorktreeCount int64
}

// Snapshot is a plain collect result (not Prometheus).
type Snapshot struct {
	ResolvedRoot    string
	Present         bool
	TotalBytes      int64
	Layouts         map[string]LayoutStats // keys: static-cache, worktrees, cache, targeted-cache
	TmpWorktrees    LayoutStats
	TmpResolvedRoot string
	ScanDuration    time.Duration
	Incomplete      bool
	Err             error
}

// layoutKeys is the ordered set of unified layout segments under ReposRoot.
var layoutKeys = []string{
	LayoutStaticCache,
	LayoutWorktrees,
	LayoutCache,
	LayoutTargetedCache,
}

// scanMode selects how directories under a layout root are classified.
type scanMode int

const (
	// modeBare counts bare git repositories (static-cache / cache / targeted-cache).
	modeBare scanMode = iota
	// modeWorktree counts worktree checkouts (.git file or directory).
	modeWorktree
)

// Collect walks injectable roots and returns best-effort inventory.
// Prefer opts-only (no context); soft budget via MaxWalkDuration.
func Collect(opts Options) Snapshot {
	start := time.Now()
	snap := Snapshot{
		Layouts: make(map[string]LayoutStats, len(layoutKeys)),
	}

	root, err := resolveReposRoot(opts.ReposRoot)
	if err != nil {
		snap.Err = err
		snap.ScanDuration = time.Since(start)
		// Still record empty layout keys for stable map shape.
		for _, k := range layoutKeys {
			snap.Layouts[k] = LayoutStats{}
		}
		return snap
	}
	root = filepath.Clean(root)
	snap.ResolvedRoot = root

	budget := &walkBudget{
		start:       start,
		maxDuration: opts.MaxWalkDuration,
		incomplete:  &snap.Incomplete,
	}

	// Root presence: missing is not an error.
	if fi, err := os.Stat(root); err != nil || !fi.IsDir() {
		snap.Present = false
		for _, k := range layoutKeys {
			snap.Layouts[k] = LayoutStats{}
		}
		// Optionally still resolve/record tmp root path when requested.
		if opts.IncludeTmpWorktrees {
			tmpRoot := resolveTmpRoot(opts.TmpWorktreesRoot)
			snap.TmpResolvedRoot = filepath.Clean(tmpRoot)
			if !budget.exceeded() {
				snap.TmpWorktrees = scanLayout(tmpRoot, modeWorktree, budget)
			}
		}
		snap.ScanDuration = time.Since(start)
		return snap
	}
	snap.Present = true
	// Soft budget may already be exhausted after the root stat (e.g. MaxWalkDuration=1ns).
	// Calling exceeded here ensures Incomplete is sticky even if subsequent walks short-circuit.
	_ = budget.exceeded()

	for _, key := range layoutKeys {
		layoutPath := filepath.Join(root, key)
		mode := modeBare
		if key == LayoutWorktrees {
			mode = modeWorktree
		}
		if budget.exceeded() {
			// Mark Present based on existence even if we skip the walk.
			if fi, err := os.Stat(layoutPath); err == nil && fi.IsDir() {
				snap.Layouts[key] = LayoutStats{Present: true}
			} else {
				snap.Layouts[key] = LayoutStats{}
			}
			continue
		}
		st := scanLayout(layoutPath, mode, budget)
		snap.Layouts[key] = st
		snap.TotalBytes += st.Bytes
	}

	if opts.IncludeTmpWorktrees {
		tmpRoot := resolveTmpRoot(opts.TmpWorktreesRoot)
		snap.TmpResolvedRoot = filepath.Clean(tmpRoot)
		if !budget.exceeded() {
			snap.TmpWorktrees = scanLayout(tmpRoot, modeWorktree, budget)
		} else if fi, err := os.Stat(tmpRoot); err == nil && fi.IsDir() {
			// Soft-timeout after unified layouts: record presence without counts.
			snap.TmpWorktrees = LayoutStats{Present: true}
		}
	}

	snap.ScanDuration = time.Since(start)
	return snap
}

func resolveReposRoot(inject string) (string, error) {
	if inject != "" {
		return inject, nil
	}
	return repos.ReposRoot()
}

func resolveTmpRoot(inject string) string {
	if inject != "" {
		return inject
	}
	return repos.TmpWorktreesRoot()
}

type walkBudget struct {
	start       time.Time
	maxDuration time.Duration // 0 = unlimited
	incomplete  *bool
}

func (b *walkBudget) exceeded() bool {
	if b == nil || b.maxDuration <= 0 {
		return false
	}
	// Inclusive: MaxWalkDuration=1ns must trip as soon as any real time elapses.
	if time.Since(b.start) >= b.maxDuration {
		if b.incomplete != nil {
			*b.incomplete = true
		}
		return true
	}
	return false
}

// scanLayout walks layoutPath and accumulates du-style bytes plus repo/worktree counts.
// Missing path → Present=false, zeros. Existing dir → Present=true even if empty.
func scanLayout(layoutPath string, mode scanMode, budget *walkBudget) LayoutStats {
	fi, err := os.Stat(layoutPath)
	if err != nil || !fi.IsDir() {
		return LayoutStats{}
	}
	st := LayoutStats{Present: true}
	if budget.exceeded() {
		return st
	}
	walkDir(layoutPath, mode, &st, budget)
	return st
}

// walkDir recursively scans path. On bare/worktree hits it still rolls up
// subtree bytes but does not look for further classified repos inside them.
func walkDir(path string, mode scanMode, st *LayoutStats, budget *walkBudget) {
	if budget.exceeded() {
		return
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}
	for _, e := range entries {
		if budget.exceeded() {
			return
		}
		full := filepath.Join(path, e.Name())
		if e.IsDir() {
			switch mode {
			case modeBare:
				if isBareRepository(full) {
					st.RepoCount++
					// Count all bytes under the bare; do not look for nested bares.
					addBytesRecursive(full, st, budget)
					continue
				}
			case modeWorktree:
				if isWorktreeCheckout(full) {
					st.WorktreeCount++
					// Count checkout bytes; skip nested classification.
					addBytesRecursive(full, st, budget)
					continue
				}
			}
			walkDir(full, mode, st, budget)
			continue
		}
		// Regular file (or symlink etc.): add size when available.
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.Mode().IsRegular() {
			st.Bytes += info.Size()
		}
	}
}

// addBytesRecursive sums regular-file sizes under root without classifying repos.
// Uses a manual walk (no filepath.SkipAll) for go1.19 compatibility.
func addBytesRecursive(root string, st *LayoutStats, budget *walkBudget) {
	if budget.exceeded() {
		return
	}
	var walk func(string)
	walk = func(dir string) {
		if budget.exceeded() {
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			if budget.exceeded() {
				return
			}
			full := filepath.Join(dir, e.Name())
			if e.IsDir() {
				walk(full)
				continue
			}
			info, err := e.Info()
			if err != nil {
				continue
			}
			if info.Mode().IsRegular() {
				st.Bytes += info.Size()
			}
		}
	}
	walk(root)
}

// isBareRepository reports whether dir is a bare git repository.
// Reimplements the unexported repos.isBareRepository check via git rev-parse.
// Fast-path: require a HEAD entry so non-git dirs skip the subprocess.
func isBareRepository(dir string) bool {
	if _, err := os.Lstat(filepath.Join(dir, "HEAD")); err != nil {
		return false
	}
	cmd := exec.Command("git", "-c", "core.hooksPath=/dev/null", "rev-parse", "--is-bare-repository")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

// isWorktreeCheckout reports a leaf checkout with a .git file or .git directory.
func isWorktreeCheckout(dir string) bool {
	gitPath := filepath.Join(dir, ".git")
	info, err := os.Lstat(gitPath)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular() || info.IsDir()
}
