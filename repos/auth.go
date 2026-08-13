package repos

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

// GitAuthConfig carries HTTPS authentication settings for git operations.
type GitAuthConfig struct {
	Token       string
	ExtraHeader string
}

type gitURLInfo struct {
	Token    string
	CleanURL string
	HasToken bool
}

func extractTokenFromURL(url string) (*gitURLInfo, error) {
	info := &gitURLInfo{HasToken: false}
	pattern := `^(https?://)(?:([^:]+):)?([^@]+)@(.+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(url)
	if matches == nil {
		info.CleanURL = url
		return info, nil
	}
	info.Token = matches[3]
	info.HasToken = true
	info.CleanURL = matches[1] + matches[4]
	return info, nil
}

func generateAuthHeader(token string) string {
	if token == "" {
		return ""
	}
	credentials := fmt.Sprintf("oauth2:%s", token)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return fmt.Sprintf("http.extraHeader=Authorization: Basic %s", encoded)
}

func maskSensitiveInfo(message, token string) string {
	if token == "" {
		return message
	}
	masked := message
	if len(token) > 8 {
		maskedToken := token[:4] + "****" + token[len(token)-4:]
		masked = strings.ReplaceAll(masked, token, maskedToken)
	} else if len(token) > 0 {
		masked = strings.ReplaceAll(masked, token, fmt.Sprintf("****(%d chars)", len(token)))
	}
	authHeader := generateAuthHeader(token)
	if authHeader != "" {
		masked = strings.ReplaceAll(masked, authHeader, "http.extraHeader=Authorization: Basic ****")
	}
	return masked
}

// AuthFromCloneURL extracts auth settings and a clean URL from a clone URL.
func AuthFromCloneURL(cloneURL string) (*GitAuthConfig, string, error) {
	return authFromCloneURL(cloneURL)
}

func authFromCloneURL(cloneURL string) (*GitAuthConfig, string, error) {
	info, err := extractTokenFromURL(cloneURL)
	if err != nil {
		return nil, cloneURL, err
	}
	if !info.HasToken {
		return nil, cloneURL, nil
	}
	return &GitAuthConfig{
		Token:       info.Token,
		ExtraHeader: generateAuthHeader(info.Token),
	}, info.CleanURL, nil
}

func (a *GitAuthConfig) gitArgs(args ...string) []string {
	if a == nil || a.ExtraHeader == "" {
		return args
	}
	withAuth := make([]string, 0, len(args)+2)
	withAuth = append(withAuth, "-c", a.ExtraHeader)
	withAuth = append(withAuth, args...)
	return withAuth
}

func (a *GitAuthConfig) maskError(err error) string {
	if err == nil {
		return ""
	}
	if a == nil || a.Token == "" {
		return err.Error()
	}
	return maskSensitiveInfo(err.Error(), a.Token)
}

func mergeAuth(optAuth *GitAuthConfig, cloneURL string) (*GitAuthConfig, string, error) {
	if optAuth != nil {
		return optAuth, cloneURL, nil
	}
	return authFromCloneURL(cloneURL)
}

// BuildAuthURL injects an OAuth2 token into a git HTTPS URL.
func BuildAuthURL(repoURL, token string) string {
	if token == "" {
		return repoURL
	}
	re := regexp.MustCompile(`(https?://)([^/]+)(/.+)$`)
	m := re.FindStringSubmatch(strings.TrimSpace(repoURL))
	if m == nil {
		return repoURL
	}
	return m[1] + "oauth2:" + token + "@" + m[2] + m[3]
}

// AuthFromToken builds git auth config and an authenticated clone URL.
func AuthFromToken(repoURL, token string) (*GitAuthConfig, string) {
	authURL := BuildAuthURL(repoURL, token)
	auth, _, err := authFromCloneURL(authURL)
	if err != nil || auth == nil {
		return nil, authURL
	}
	return auth, authURL
}

func cloneURLForGit(cloneURL string, auth *GitAuthConfig) string {
	info, err := extractTokenFromURL(cloneURL)
	if err == nil && info.HasToken {
		return info.CleanURL
	}
	return cloneURL
}

func isBareRepository(dir string) (bool, error) {
	out, err := gitOutput(dir, "rev-parse", "--is-bare-repository")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) == "true", nil
}