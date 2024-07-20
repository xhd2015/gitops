package url

import (
	"fmt"
	"regexp"
	"strings"
)

var httpRegex = regexp.MustCompile(`^https?://(?:[^:]+:[^@]+@)?([^/]+(:\d+)?)(/.+?)(\.git)?$`)
var sshRegex = regexp.MustCompile(`^[^@]+@([^:]+):(.+?)(\.git)?$`)

// SplitRepoURL split repo into two parts: domain, path, without 'https://', 'http://', or 'ssh://'
// example:  https://gitlab:TOKEN@xxxx
//
//	gitlab@xxx.com:/xx/xxx.git
func SplitRepoURL(repoURL string) (domain string, path string, err error) {
	if repoURL == "" {
		err = fmt.Errorf("empty repoURL")
		return
	}
	if r := httpRegex.FindStringSubmatch(repoURL); r != nil {
		// fmt.Printf("r: %+v", strings.Join(r, ";"))
		domain = r[1]
		path = r[3]
		return
	}
	if r := sshRegex.FindStringSubmatch(repoURL); r != nil {
		// fmt.Printf("r2: %+v", strings.Join(r, ";"))
		domain = r[1]
		path = r[2]
		if path != "" && !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		return
	}

	err = fmt.Errorf("unrecognized repoURL")
	return
}

func JoinRepoURL(repoURL string, user string, token string) (string, error) {
	if repoURL == "" {
		return "", fmt.Errorf("empty repoURL")
	}
	domain, path, err := SplitRepoURL(repoURL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s:%s@%s/%s", user, token, domain, strings.TrimPrefix(path, "/")), nil
}