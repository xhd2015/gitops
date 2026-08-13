package repos

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	// DefaultTargetRefPrefix is the dest-ref namespace when TargetRefPrefix is empty.
	DefaultTargetRefPrefix     = "refs/gitops/targets"
	defaultTargetFreshness     = time.Minute
	targetedRemoteCheckTimeout = 5 * time.Second
)

var fullCommitPattern = regexp.MustCompile(`^[0-9a-fA-F]{40}$`)

// ErrTargetRefNotFound identifies an advertised branch/tag that does not
// exist. Callers may preserve an exists=false API response for this case.
var ErrTargetRefNotFound = errors.New("target git ref not found")

// IsTargetRefNotFound reports whether a targeted fetch failed because its
// requested branch/tag is absent.
func IsTargetRefNotFound(err error) bool {
	return errors.Is(err, ErrTargetRefNotFound)
}

// TargetedFetchOptions controls a narrow fetch into a partial bare cache.
type TargetedFetchOptions struct {
	Auth            *GitAuthConfig
	Env             map[string]string
	MaxAge          time.Duration
	Depth           int
	Deepen          int
	Unshallow       bool
	HistoryPhase    string
	ReposRoot       string    // empty → ReposRoot(); cache under targeted-cache/
	Progress        io.Writer // dir: + $ git … before fetch; have <sha> when present
	TargetRefPrefix string    // empty → DefaultTargetRefPrefix; trailing / trimmed
}

// TargetedFetchResult describes the commits made available by a narrow fetch.
type TargetedFetchResult struct {
	CacheDir       string
	Commits        map[string]string
	Fetched        bool
	FilterAccepted bool
	LockWait       time.Duration
	FetchDuration  time.Duration
}

type targetedRef struct {
	input     string
	source    string
	local     string
	immutable bool
}

// EnsureTargetedCache initializes or reuses an isolated promisor bare cache.
func EnsureTargetedCache(ctx context.Context, cloneURL string, opts FetchOptions) (string, error) {
	cacheDir, err := TargetedCacheDirUnder(opts.ReposRoot, cloneURL)
	if err != nil {
		return "", err
	}
	if targetedCacheReady(cacheDir) {
		return cacheDir, nil
	}

	auth, cleanURL, err := mergeAuth(opts.Auth, cloneURL)
	if err != nil {
		return "", fmt.Errorf("parse auth: %w", err)
	}
	lock, err := acquireRepoLockContext(ctx, cacheDir)
	if err != nil {
		return "", err
	}
	defer releaseRepoLock(lock)
	if ok, _ := isBareRepository(cacheDir); ok {
		if err := configureTargetedCache(ctx, cacheDir, cloneURLForGit(cleanURL, auth), opts.Env); err != nil {
			return "", err
		}
		return cacheDir, nil
	}
	if _, err := os.Stat(cacheDir); err == nil {
		return "", fmt.Errorf("targeted cache path exists but is not a bare repository: %s", cacheDir)
	} else if !os.IsNotExist(err) {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(cacheDir), 0o755); err != nil {
		return "", err
	}
	if err := runGitWithEnvContext(ctx, "", auth, opts.Env, "init", "--bare", cacheDir); err != nil {
		return "", fmt.Errorf("init targeted cache: %w", err)
	}
	if err := configureTargetedCache(ctx, cacheDir, cloneURLForGit(cleanURL, auth), opts.Env); err != nil {
		return "", err
	}
	return cacheDir, nil
}

func targetedCacheReady(cacheDir string) bool {
	if ok, _ := isBareRepository(cacheDir); !ok {
		return false
	}
	promisor, err := gitOutput(cacheDir, "config", "--bool", "--get", "remote.origin.promisor")
	if err != nil || strings.TrimSpace(promisor) != "true" {
		return false
	}
	_, err = gitOutput(cacheDir, "remote", "get-url", "origin")
	return err == nil
}

func configureTargetedCache(ctx context.Context, cacheDir, cleanURL string, env map[string]string) error {
	if _, err := gitOutput(cacheDir, "remote", "get-url", "origin"); err == nil {
		if err := runGitWithEnvContext(ctx, cacheDir, nil, env, "remote", "set-url", "origin", cleanURL); err != nil {
			return fmt.Errorf("set targeted cache origin: %w", err)
		}
	} else if err := runGitWithEnvContext(ctx, cacheDir, nil, env, "remote", "add", "origin", cleanURL); err != nil {
		return fmt.Errorf("add targeted cache origin: %w", err)
	}
	// git remote add installs +refs/heads/*:refs/remotes/origin/*. Explicit dest
	// refspecs still trigger opportunistic updates of those tracking refs; drop
	// the default so a narrow fetch cannot create refs/remotes/origin/*.
	if err := runGitWithEnvContext(ctx, cacheDir, nil, env, "config", "--unset-all", "remote.origin.fetch"); err != nil && !strings.Contains(err.Error(), "exit status 5") {
		return fmt.Errorf("configure targeted cache: %w", err)
	}
	configure := [][]string{
		{"config", "remote.origin.promisor", "true"},
		{"config", "remote.origin.partialclonefilter", "blob:none"},
		{"config", "extensions.partialClone", "origin"},
	}
	for _, args := range configure {
		if err := runGitWithEnvContext(ctx, cacheDir, nil, env, args...); err != nil {
			return fmt.Errorf("configure targeted cache: %w", err)
		}
	}
	return nil
}

// FetchTargetedRefs ensures only the requested commits/branches/tags and their
// required history are present in the isolated partial cache.
func FetchTargetedRefs(ctx context.Context, cloneURL string, refs []string, opts TargetedFetchOptions) (*TargetedFetchResult, error) {
	auth, _, err := mergeAuth(opts.Auth, cloneURL)
	if err != nil {
		return nil, err
	}
	cacheDir, err := EnsureTargetedCache(ctx, cloneURL, FetchOptions{Auth: auth, Env: opts.Env, ReposRoot: opts.ReposRoot})
	if err != nil {
		return nil, err
	}
	targets, err := normalizeTargetedRefs(refs, opts.TargetRefPrefix)
	if err != nil {
		return nil, err
	}
	result := &TargetedFetchResult{
		CacheDir:       cacheDir,
		Commits:        make(map[string]string, len(targets)),
		FilterAccepted: true,
	}

	if !opts.Unshallow && opts.Deepen == 0 {
		fresh, err := targetedRefsFresh(cacheDir, targets, opts.MaxAge)
		if err == nil && fresh {
			if err := resolveTargetedCommits(cacheDir, targets, result.Commits); err == nil {
				printHaveTargets(opts.Progress, targets)
				return result, nil
			}
		}
		if err := resolveTargetedCommits(cacheDir, targets, result.Commits); err == nil {
			unchanged, err := targetedRemoteRefsUnchanged(ctx, cacheDir, auth, opts.Env, targets, result.Commits)
			if err != nil {
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				if !targetedRemoteRefreshUnavailable(err) {
					return nil, err
				}
				// The requested commits are already present and the caller was
				// authorized to obtain a clone URL. Prefer a bounded stale
				// result during transient Git transport failures instead of
				// blocking every read on the same unavailable remote.
				if err := writeTargetedFreshness(cacheDir, targets); err != nil {
					return nil, err
				}
				printHaveTargets(opts.Progress, targets)
				return result, nil
			}
			if unchanged {
				if err := writeTargetedFreshness(cacheDir, targets); err != nil {
					return nil, err
				}
				printHaveTargets(opts.Progress, targets)
				return result, nil
			}
		}
	}

	lockStart := time.Now()
	lock, err := acquireRepoLockContext(ctx, cacheDir)
	if err != nil {
		return nil, err
	}
	defer releaseRepoLock(lock)
	result.LockWait = time.Since(lockStart)

	if !opts.Unshallow && opts.Deepen == 0 {
		fresh, err := targetedRefsFresh(cacheDir, targets, opts.MaxAge)
		if err == nil && fresh {
			if err := resolveTargetedCommits(cacheDir, targets, result.Commits); err == nil {
				printHaveTargets(opts.Progress, targets)
				return result, nil
			}
		}
	}
	if opts.HistoryPhase != "" {
		if len(targets) != 2 {
			return nil, fmt.Errorf("history connectivity recheck requires exactly two refs")
		}
		if err := resolveTargetedCommits(cacheDir, targets, result.Commits); err == nil {
			connected, err := TargetedHistoryConnected(
				ctx,
				cacheDir,
				result.Commits[targets[0].input],
				result.Commits[targets[1].input],
			)
			if err != nil {
				return nil, err
			}
			if connected {
				return result, nil
			}
			phaseDone, err := targetedHistoryPhaseDone(cacheDir, opts.HistoryPhase, result.Commits, targets)
			if err != nil {
				return nil, err
			}
			if phaseDone {
				return result, nil
			}
			if opts.Unshallow {
				shallow, err := gitOutput(cacheDir, "rev-parse", "--is-shallow-repository")
				if err != nil {
					return nil, err
				}
				// Another caller may have completed the final unshallow while
				// this caller waited. For unrelated histories there is still
				// no merge base, but another full fetch cannot change that.
				if strings.TrimSpace(shallow) == "false" {
					return result, nil
				}
			}
		}
	}

	args := []string{"fetch", "--no-tags", "--filter=blob:none"}
	switch {
	case opts.Unshallow:
		shallow, _ := gitOutput(cacheDir, "rev-parse", "--is-shallow-repository")
		if strings.TrimSpace(shallow) == "true" {
			args = append(args, "--unshallow")
		}
	case opts.Deepen > 0:
		args = append(args, fmt.Sprintf("--deepen=%d", opts.Deepen))
	default:
		depth := opts.Depth
		if depth <= 0 {
			depth = 1
		}
		args = append(args, fmt.Sprintf("--depth=%d", depth))
	}
	args = append(args, "origin")
	for _, target := range targets {
		args = append(args, "+"+target.source+":"+target.local)
	}
	fetchStart := time.Now()
	out, err := runTargetedFetch(ctx, cacheDir, auth, opts.Env, opts.Progress, args)
	result.FetchDuration = time.Since(fetchStart)
	if err != nil {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "couldn't find remote ref") ||
			strings.Contains(msg, "remote ref does not exist") ||
			strings.Contains(msg, "not our ref") ||
			strings.Contains(msg, "unadvertised object") {
			return nil, fmt.Errorf("%w: %s", ErrTargetRefNotFound, auth.maskError(err))
		}
		return nil, fmt.Errorf("fetch targeted refs: %s", auth.maskError(err))
	}
	result.Fetched = true
	if strings.Contains(strings.ToLower(out), "filtering not recognized") {
		result.FilterAccepted = false
	}
	if err := writeTargetedFreshness(cacheDir, targets); err != nil {
		return nil, err
	}
	if err := resolveTargetedCommits(cacheDir, targets, result.Commits); err != nil {
		return nil, err
	}
	if opts.HistoryPhase != "" {
		if err := markTargetedHistoryPhaseDone(cacheDir, opts.HistoryPhase, result.Commits, targets); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// ReadTargetedFile reads one file from a commit in a promisor cache. Holding
// the mutation lock also serializes Git's implicit lazy blob fetch.
func ReadTargetedFile(ctx context.Context, cacheDir string, auth *GitAuthConfig, env map[string]string, commit, file string) (bool, string, error) {
	if commit == "" {
		return false, "", fmt.Errorf("requires commit")
	}
	if file == "" {
		return false, "", fmt.Errorf("requires file")
	}
	lock, err := acquireRepoLockContext(ctx, cacheDir)
	if err != nil {
		return false, "", err
	}
	defer releaseRepoLock(lock)
	spec := commit + ":" + file
	out, err := gitStdoutWithEnvContext(ctx, cacheDir, auth, env, "cat-file", "blob", spec)
	if err != nil {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "does not exist") ||
			strings.Contains(msg, "not a valid object name") ||
			strings.Contains(msg, "path '"+strings.ToLower(file)+"' does not exist") {
			return false, "", nil
		}
		return false, "", fmt.Errorf("read targeted file: %s", auth.maskError(err))
	}
	return true, string(out), nil
}

// TargetedCommitAvailable reports whether commit is already reachable from
// the refs held in a partial cache. --missing=allow-any prevents this probe
// from triggering Git's implicit promisor fetch for an absent object.
func TargetedCommitAvailable(ctx context.Context, cacheDir, commit string) (bool, error) {
	if commit == "" {
		return false, fmt.Errorf("requires commit")
	}
	out, err := gitOutputWithEnvContext(ctx, cacheDir, nil, nil, "rev-list", "--all", "--missing=allow-any")
	if err != nil {
		return false, fmt.Errorf("list targeted commits: %w", err)
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == commit {
			return true, nil
		}
	}
	return false, nil
}

// TargetedHistoryConnected reports whether the partial cache contains a
// common ancestor for the two commits. A missing merge base in a shallow
// repository means the caller must deepen before trusting graph results.
func TargetedHistoryConnected(ctx context.Context, cacheDir, left, right string) (bool, error) {
	_, err := gitOutputWithEnvContext(ctx, cacheDir, nil, nil, "merge-base", left, right)
	if err == nil {
		return true, nil
	}
	if strings.Contains(err.Error(), "exit status 1") {
		return false, nil
	}
	return false, err
}

// NormalizeTargetRefPrefix returns DefaultTargetRefPrefix when prefix is empty
// and strips trailing slashes so dest refs never contain "//".
func NormalizeTargetRefPrefix(prefix string) string {
	p := strings.TrimRight(prefix, "/")
	if p == "" {
		return DefaultTargetRefPrefix
	}
	return p
}

// TargetDestRef is the dest local ref: {normalizedPrefix}/{hex(sha256(input))}.
func TargetDestRef(prefix, input string) string {
	sum := sha256.Sum256([]byte(input))
	return NormalizeTargetRefPrefix(prefix) + "/" + hex.EncodeToString(sum[:])
}

func normalizeTargetedRefs(refs []string, prefix string) ([]targetedRef, error) {
	prefix = NormalizeTargetRefPrefix(prefix)
	seen := make(map[string]bool, len(refs))
	targets := make([]targetedRef, 0, len(refs))
	for _, input := range refs {
		if input == "" || seen[input] {
			continue
		}
		seen[input] = true
		source := input
		immutable := fullCommitPattern.MatchString(input)
		switch {
		case immutable:
		case strings.HasPrefix(input, "refs/heads/"):
		case strings.HasPrefix(input, "refs/tags/"):
		case strings.HasPrefix(input, "refs/remotes/origin/"):
			source = "refs/heads/" + strings.TrimPrefix(input, "refs/remotes/origin/")
		case strings.HasPrefix(input, "origin/"):
			source = "refs/heads/" + strings.TrimPrefix(input, "origin/")
		case !strings.HasPrefix(input, "refs/"):
			source = "refs/heads/" + input
		default:
			return nil, fmt.Errorf("unsupported git ref %q", input)
		}
		if !immutable {
			if err := runGit("", nil, "check-ref-format", source); err != nil {
				return nil, fmt.Errorf("invalid git ref %q", input)
			}
		}
		targets = append(targets, targetedRef{
			input:     input,
			source:    source,
			local:     TargetDestRef(prefix, input),
			immutable: immutable,
		})
	}
	sort.Slice(targets, func(i, j int) bool { return targets[i].input < targets[j].input })
	if len(targets) == 0 {
		return nil, fmt.Errorf("requires at least one git ref")
	}
	return targets, nil
}

func targetedRefsFresh(cacheDir string, targets []targetedRef, maxAge time.Duration) (bool, error) {
	if maxAge == 0 {
		maxAge = defaultTargetFreshness
	}
	for _, target := range targets {
		if _, err := gitOutput(cacheDir, "rev-parse", "--verify", target.local+"^{commit}"); err != nil {
			return false, nil
		}
		if target.immutable {
			continue
		}
		if maxAge < 0 {
			return false, nil
		}
		info, err := os.Stat(targetedFreshnessPath(cacheDir, target))
		if err != nil || time.Since(info.ModTime()) >= maxAge {
			return false, nil
		}
	}
	return true, nil
}

func writeTargetedFreshness(cacheDir string, targets []targetedRef) error {
	dir := filepath.Join(cacheDir, "gitops-target-state")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	now := time.Now()
	for _, target := range targets {
		path := targetedFreshnessPath(cacheDir, target)
		if err := os.WriteFile(path, []byte(target.input+"\n"), 0o644); err != nil {
			return err
		}
		if err := os.Chtimes(path, now, now); err != nil {
			return err
		}
	}
	return nil
}

func targetedRemoteRefsUnchanged(
	ctx context.Context,
	cacheDir string,
	auth *GitAuthConfig,
	env map[string]string,
	targets []targetedRef,
	commits map[string]string,
) (bool, error) {
	args := []string{"ls-remote", "origin"}
	for _, target := range targets {
		if target.immutable {
			continue
		}
		args = append(args, target.source)
	}
	if len(args) == 2 {
		return true, nil
	}
	checkCtx, cancel := context.WithTimeout(ctx, targetedRemoteCheckTimeout)
	defer cancel()
	out, err := gitOutputWithEnvContext(checkCtx, cacheDir, auth, env, args...)
	if err != nil {
		if checkCtx.Err() != nil && ctx.Err() == nil {
			return false, fmt.Errorf("check targeted remote refs: %w", checkCtx.Err())
		}
		return false, fmt.Errorf("check targeted remote refs: %s", auth.maskError(err))
	}
	advertised := make(map[string]string, len(targets))
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		advertised[fields[1]] = fields[0]
	}
	for _, target := range targets {
		if target.immutable {
			continue
		}
		remoteCommit := advertised[target.source]
		if peeled := advertised[target.source+"^{}"]; peeled != "" {
			remoteCommit = peeled
		}
		if remoteCommit == "" || remoteCommit != commits[target.input] {
			return false, nil
		}
	}
	return true, nil
}

func targetedRemoteRefreshUnavailable(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	message := strings.ToLower(err.Error())
	for _, fragment := range []string{
		"could not resolve host",
		"failed to connect",
		"connection timed out",
		"operation timed out",
		"recv failure",
		"connection reset",
		"network is unreachable",
	} {
		if strings.Contains(message, fragment) {
			return true
		}
	}
	return false
}

func targetedFreshnessPath(cacheDir string, target targetedRef) string {
	return filepath.Join(cacheDir, "gitops-target-state", filepath.Base(target.local))
}

func targetedHistoryPhaseDone(cacheDir, phase string, commits map[string]string, targets []targetedRef) (bool, error) {
	path, err := targetedHistoryPhasePath(cacheDir, phase, commits, targets)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func markTargetedHistoryPhaseDone(cacheDir, phase string, commits map[string]string, targets []targetedRef) error {
	path, err := targetedHistoryPhasePath(cacheDir, phase, commits, targets)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(phase+"\n"), 0o644)
}

func targetedHistoryPhasePath(cacheDir, phase string, commits map[string]string, targets []targetedRef) (string, error) {
	if phase == "" || len(targets) != 2 {
		return "", fmt.Errorf("history phase requires a name and exactly two refs")
	}
	pair := []string{commits[targets[0].input], commits[targets[1].input]}
	if pair[0] == "" || pair[1] == "" {
		return "", fmt.Errorf("history phase requires two resolved commits")
	}
	sort.Strings(pair)
	sum := sha256.Sum256([]byte(phase + "\x00" + pair[0] + "\x00" + pair[1]))
	return filepath.Join(cacheDir, "gitops-history-state", hex.EncodeToString(sum[:])), nil
}

func printHaveTargets(w io.Writer, targets []targetedRef) {
	if w == nil {
		return
	}
	for _, target := range targets {
		fmt.Fprintf(w, "have %s\n", target.input)
	}
}

// runTargetedFetch runs git fetch. When Progress is set, prints dir: and the
// exact $ git … (with --progress) first, then streams via runGitStreamStderr.
func runTargetedFetch(ctx context.Context, cacheDir string, auth *GitAuthConfig, env map[string]string, progress io.Writer, args []string) (string, error) {
	if progress == nil {
		return gitOutputWithEnvContext(ctx, cacheDir, auth, env, args...)
	}
	runArgs := make([]string, 0, len(args)+1)
	if len(args) > 0 {
		runArgs = append(runArgs, args[0], "--progress")
		runArgs = append(runArgs, args[1:]...)
	} else {
		runArgs = []string{"fetch", "--progress"}
	}
	fmt.Fprintf(progress, "dir: %s\n", cacheDir)
	fmt.Fprintf(progress, "$ git %s\n", strings.Join(runArgs, " "))
	var captured bytes.Buffer
	w := io.MultiWriter(progress, &captured)
	err := runGitStreamStderr(ctx, cacheDir, auth, env, w, runArgs...)
	return captured.String(), err
}

func resolveTargetedCommits(cacheDir string, targets []targetedRef, dest map[string]string) error {
	for _, target := range targets {
		commit, err := gitOutput(cacheDir, "rev-parse", "--verify", target.local+"^{commit}")
		if err != nil {
			return fmt.Errorf("resolve targeted ref %s: %w", target.input, err)
		}
		dest[target.input] = strings.TrimSpace(commit)
	}
	return nil
}
