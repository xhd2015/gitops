package git

import (
	"fmt"
	"strings"

	"github.com/xhd2015/xgo/support/cmd"
)

const REFS_REMOTES_ORIGIN_PREFIX = "refs/remotes/origin/"

// see https://stackoverflow.com/questions/2706797/finding-what-branch-a-git-commit-came-from
func GetBranchesContainingRef(dir string, ref string) (branches []string, err error) {
	if dir == "" {
		return nil, fmt.Errorf("requires dir")
	}
	if ref == "" {
		return nil, fmt.Errorf("requires ref")
	}
	// example:
	//   $ git branch -a --contains release-v2.16.0 --format='%(refname)'
	//    refs/heads/release-v2.16.0
	//    refs/heads/release-v2.18.0
	//    refs/heads/release-v2.21.0
	//    refs/remotes/origin/HEAD
	//    refs/remotes/origin/master
	//    refs/remotes/origin/release-v2.16.0
	//    refs/remotes/origin/release-v2.17.1
	//    refs/remotes/origin/release-v2.18.0
	//    refs/remotes/origin/release-v2.21.0
	output, err := cmd.Dir(dir).Output("git", "branch", "-a", "--format=%(refname)", "--contains", ref)
	if err != nil {
		return nil, err
	}

	var candidates []string
	lines := splitLinesFilterEmpty(output)
	seen := make(map[string]bool, len(lines))
	for _, line := range lines {
		if !strings.HasPrefix(line, REFS_REMOTES_ORIGIN_PREFIX) {
			continue
		}
		branch := line[len(REFS_REMOTES_ORIGIN_PREFIX):]
		if branch == "HEAD" {
			continue
		}
		if seen[branch] {
			continue
		}
		seen[branch] = true
		candidates = append(candidates, branch)
	}

	for _, candidate := range candidates {
		isAncestor, err := IsFirstParentAncestorOf(dir, REFS_REMOTES_ORIGIN_PREFIX+candidate, ref)
		if err != nil {
			return nil, err
		}
		if !isAncestor {
			continue
		}
		branches = append(branches, candidate)
	}
	return branches, nil
}