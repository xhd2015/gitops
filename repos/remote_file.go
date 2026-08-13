package repos

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const maxRemoteFileBytes = 32 << 20

// ReadGitLabFile reads one file directly from the GitLab API at an exact
// commit/ref. It avoids creating a repository cache or taking a repository
// lock. Authentication uses PRIVATE-TOKEN (the token extracted from the clone
// URL / auth config), not Basic oauth2 against the web UI raw path.
func ReadGitLabFile(ctx context.Context, cloneURL, commit, file string, opts FetchOptions) (bool, string, error) {
	if commit == "" {
		return false, "", fmt.Errorf("requires commit")
	}
	file = strings.TrimPrefix(file, "/")
	if file == "" {
		return false, "", fmt.Errorf("requires file")
	}
	for _, part := range strings.Split(file, "/") {
		if part == "" || part == "." || part == ".." {
			return false, "", fmt.Errorf("invalid file path")
		}
	}

	auth, cleanURL, err := mergeAuth(opts.Auth, cloneURL)
	if err != nil {
		return false, "", err
	}
	// Prefer the clean clone URL (no embedded credentials) for path parsing.
	if auth != nil {
		if _, stripped, authErr := authFromCloneURL(cloneURL); authErr == nil && stripped != "" {
			cleanURL = stripped
		}
	}
	repoURL, err := url.Parse(cleanURL)
	if err != nil {
		return false, "", fmt.Errorf("parse clone URL: %w", err)
	}
	if repoURL.Scheme != "http" && repoURL.Scheme != "https" {
		return false, "", fmt.Errorf("GitLab raw file requires an HTTP(S) clone URL")
	}

	projectPath := strings.TrimPrefix(repoURL.EscapedPath(), "/")
	if projectPath == "" {
		projectPath = strings.TrimPrefix(repoURL.Path, "/")
	}
	projectPath = strings.TrimSuffix(projectPath, ".git")
	projectPath = strings.Trim(projectPath, "/")
	if projectPath == "" {
		return false, "", fmt.Errorf("GitLab project path is empty")
	}
	// EscapedPath may still contain unescaped segments when Path was set only;
	// normalize to the decoded project path then re-encode as one API id.
	if unescaped, err := url.PathUnescape(projectPath); err == nil {
		projectPath = unescaped
	}

	// Build the URL as a pre-encoded string. url.URL.String() would re-escape
	// PathEscape's %2F markers and break GitLab's project/file id encoding.
	apiURL := fmt.Sprintf(
		"%s://%s/api/v4/projects/%s/repository/files/%s/raw?ref=%s",
		repoURL.Scheme,
		repoURL.Host,
		url.PathEscape(projectPath),
		url.PathEscape(file),
		url.QueryEscape(commit),
	)

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if proxy, err := proxyFromGitEnv(opts.Env); err != nil {
		return false, "", err
	} else if proxy != nil {
		transport.Proxy = http.ProxyURL(proxy)
	}
	client := &http.Client{
		Transport: transport,
		// Squid + concurrent getFile bursts routinely need >15s from worker pods.
		Timeout: 45 * time.Second,
		// Never follow redirects into the GitLab sign-in HTML page.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return false, "", err
	}
	if auth != nil && auth.Token != "" {
		// GitLab personal/project access tokens authenticate the REST API via
		// PRIVATE-TOKEN. Basic oauth2 works for git smart-HTTP, not this API.
		req.Header.Set("PRIVATE-TOKEN", auth.Token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("read GitLab raw file: %s", auth.maskError(err))
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// continue
	case http.StatusNotFound:
		// Missing project/commit/file — treat as not present so callers can
		// fall back or report exists=false.
		return false, "", nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return false, "", fmt.Errorf("read GitLab raw file: HTTP %s", resp.Status)
	default:
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			return false, "", fmt.Errorf("read GitLab raw file: unexpected redirect HTTP %s", resp.Status)
		}
		return false, "", fmt.Errorf("read GitLab raw file: HTTP %s", resp.Status)
	}

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.Contains(contentType, "text/html") {
		return false, "", fmt.Errorf("read GitLab raw file: refused HTML response (likely sign-in page)")
	}

	content, err := io.ReadAll(io.LimitReader(resp.Body, maxRemoteFileBytes+1))
	if err != nil {
		return false, "", fmt.Errorf("read GitLab raw file response: %w", err)
	}
	if len(content) > maxRemoteFileBytes {
		return false, "", fmt.Errorf("read GitLab raw file: file exceeds %d bytes", maxRemoteFileBytes)
	}
	if looksLikeHTML(content) {
		return false, "", fmt.Errorf("read GitLab raw file: refused HTML body (likely sign-in page)")
	}
	return true, string(content), nil
}

func looksLikeHTML(content []byte) bool {
	s := strings.TrimSpace(string(content))
	if len(s) > 256 {
		s = s[:256]
	}
	lower := strings.ToLower(s)
	return strings.HasPrefix(lower, "<!doctype html") ||
		strings.HasPrefix(lower, "<html") ||
		strings.Contains(lower, "devise-layout-html") ||
		(strings.Contains(lower, "<html") && strings.Contains(lower, "sign in"))
}

func proxyFromGitEnv(env map[string]string) (*url.URL, error) {
	for _, key := range []string{"https_proxy", "HTTPS_PROXY", "http_proxy", "HTTP_PROXY"} {
		raw := strings.TrimSpace(env[key])
		if raw == "" {
			continue
		}
		proxy, err := url.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", key, err)
		}
		return proxy, nil
	}
	return nil, nil
}
