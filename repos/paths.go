package repos

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ReposRoot returns the root directory for unified repo cache (~/.repos).
func ReposRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".repos"), nil
}

// TmpWorktreesRoot returns the root directory for ephemeral worktrees.
func TmpWorktreesRoot() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.TempDir(), "repos-tmp-worktrees")
	}
	return "/tmp/repos-tmp-worktrees"
}

// SanitizeRef replaces slash characters in git refs for use as path segments.
func SanitizeRef(ref string) string {
	return strings.ReplaceAll(ref, "/", "_")
}

// URLToPathSegments derives cache/worktree path segments from a clone URL.
func URLToPathSegments(cloneURL string) ([]string, error) {
	host, pathParts, err := parseCloneURL(cloneURL)
	if err != nil {
		return nil, err
	}
	segments := make([]string, 0, 1+len(pathParts))
	segments = append(segments, host)
	segments = append(segments, pathParts...)
	return segments, nil
}

// CacheDir returns the bare cache directory for a clone URL.
func CacheDir(cloneURL string) (string, error) {
	root, err := ReposRoot()
	if err != nil {
		return "", err
	}
	return CacheDirUnder(root, cloneURL)
}

// CacheDirUnder returns {reposRoot}/cache/<url-segments>/ for a clone URL.
// Empty reposRoot falls back to ReposRoot().
func CacheDirUnder(reposRoot, cloneURL string) (string, error) {
	root, err := resolveReposRoot(reposRoot)
	if err != nil {
		return "", err
	}
	segments, err := URLToPathSegments(cloneURL)
	if err != nil {
		return "", err
	}
	return filepath.Join(append([]string{root, "cache"}, segments...)...), nil
}

// TargetedCacheDir returns the isolated partial bare cache directory for a
// clone URL. It is separate from CacheDir so callers that require a complete
// mirror or persistent worktree keep their existing semantics.
func TargetedCacheDir(cloneURL string) (string, error) {
	return TargetedCacheDirUnder("", cloneURL)
}

// TargetedCacheDirUnder returns {reposRoot}/targeted-cache/<url-segments>/ for
// a clone URL. Empty reposRoot falls back to ReposRoot().
func TargetedCacheDirUnder(reposRoot, cloneURL string) (string, error) {
	root, err := resolveReposRoot(reposRoot)
	if err != nil {
		return "", err
	}
	segments, err := URLToPathSegments(cloneURL)
	if err != nil {
		return "", err
	}
	return filepath.Join(append([]string{root, "targeted-cache"}, segments...)...), nil
}

// StaticCacheDirUnder returns {reposRoot}/static-cache/<url-segments>/ for a
// clone URL. Empty reposRoot falls back to ReposRoot(). Static bares are
// siblings of cache/ and targeted-cache/ and must not share their paths.
func StaticCacheDirUnder(reposRoot, cloneURL string) (string, error) {
	root, err := resolveReposRoot(reposRoot)
	if err != nil {
		return "", err
	}
	segments, err := URLToPathSegments(cloneURL)
	if err != nil {
		return "", err
	}
	return filepath.Join(append([]string{root, "static-cache"}, segments...)...), nil
}

// WorktreeDir returns the persistent worktree directory for a clone URL and ref.
func WorktreeDir(cloneURL, ref string) (string, error) {
	root, err := ReposRoot()
	if err != nil {
		return "", err
	}
	return WorktreeDirUnder(root, cloneURL, ref)
}

// WorktreeDirUnder returns {reposRoot}/worktrees/<url-segments>/<sanitized-ref>/
// for a clone URL and ref. Empty reposRoot falls back to ReposRoot().
func WorktreeDirUnder(reposRoot, cloneURL, ref string) (string, error) {
	root, err := resolveReposRoot(reposRoot)
	if err != nil {
		return "", err
	}
	segments, err := URLToPathSegments(cloneURL)
	if err != nil {
		return "", err
	}
	segments = append(segments, SanitizeRef(ref))
	return filepath.Join(append([]string{root, "worktrees"}, segments...)...), nil
}

// EphemeralWorktreeDirUnder returns a Mode A unique worktree path:
// {reposRoot}/worktrees/<url-segments>/<sanitized-ref>/<uniqueID>/.
// Empty reposRoot falls back to ReposRoot().
func EphemeralWorktreeDirUnder(reposRoot, cloneURL, ref, uniqueID string) (string, error) {
	base, err := WorktreeDirUnder(reposRoot, cloneURL, ref)
	if err != nil {
		return "", err
	}
	return filepath.Join(base, uniqueID), nil
}

func resolveReposRoot(reposRoot string) (string, error) {
	if reposRoot != "" {
		return reposRoot, nil
	}
	return ReposRoot()
}

// TmpWorktreeDir returns the ephemeral worktree directory for a clone URL and id.
func TmpWorktreeDir(cloneURL, id string) (string, error) {
	segments, err := URLToPathSegments(cloneURL)
	if err != nil {
		return "", err
	}
	segments = append(segments, id)
	return filepath.Join(append([]string{TmpWorktreesRoot()}, segments...)...), nil
}

func parseCloneURL(cloneURL string) (host string, pathParts []string, err error) {
	if cloneURL == "" {
		return "", nil, fmt.Errorf("empty clone URL")
	}

	if strings.HasPrefix(cloneURL, "git@") {
		return parseSCPStyleURL(cloneURL)
	}

	u, err := url.Parse(cloneURL)
	if err != nil {
		return "", nil, fmt.Errorf("parse clone URL: %w", err)
	}

	switch strings.ToLower(u.Scheme) {
	case "file":
		path := u.Path
		if runtime.GOOS == "windows" && len(path) >= 3 && path[0] == '/' && path[2] == ':' {
			path = path[1:]
		}
		path = filepath.ToSlash(path)
		parts := splitPathSegments(path)
		if len(parts) == 0 {
			return "", nil, fmt.Errorf("file URL has no path segments: %s", cloneURL)
		}
		segments := append([]string{"file.local"}, parts...)
		stripGitSuffix(&segments[len(segments)-1])
		return segments[0], segments[1:], nil
	case "ssh", "http", "https":
		host = u.Hostname()
		if host == "" {
			return "", nil, fmt.Errorf("clone URL missing host: %s", cloneURL)
		}
		pathParts = splitPathSegments(u.EscapedPath())
		if len(pathParts) == 0 {
			return "", nil, fmt.Errorf("clone URL missing path: %s", cloneURL)
		}
		stripGitSuffix(&pathParts[len(pathParts)-1])
		return host, pathParts, nil
	default:
		return "", nil, fmt.Errorf("unsupported clone URL scheme %q", u.Scheme)
	}
}

func parseSCPStyleURL(cloneURL string) (host string, pathParts []string, err error) {
	at := strings.Index(cloneURL, "@")
	colon := strings.LastIndex(cloneURL, ":")
	if at < 0 || colon <= at {
		return "", nil, fmt.Errorf("invalid scp-style clone URL: %s", cloneURL)
	}
	host = cloneURL[at+1 : colon]
	pathParts = splitPathSegments(cloneURL[colon+1:])
	if len(pathParts) == 0 {
		return "", nil, fmt.Errorf("scp-style clone URL missing path: %s", cloneURL)
	}
	stripGitSuffix(&pathParts[len(pathParts)-1])
	return host, pathParts, nil
}

func splitPathSegments(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

func stripGitSuffix(segment *string) {
	*segment = strings.TrimSuffix(*segment, ".git")
}
