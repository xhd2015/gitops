package git

import (
	"fmt"

	"github.com/xhd2015/xgo/support/cmd"
)

// check is head is ref's ancestor,  can be found following head's first parent
// just like git merge-base --is-ancestor head ref, but excluding merge points
//
// see help for 'git rev-list': https://git-scm.com/docs/git-rev-list
// algorithm:
//
//	find first parent set of head, and exclude first parent set of ref^1
//	    S=git rev-list --first-parent head ^ref^1
//	if ref in S, then head is a direct parent of ref
func IsFirstParentAncestorOf(dir string, head string, ref string) (bool, error) {
	if dir == "" {
		return false, fmt.Errorf("requires dir")
	}
	refCommit, err := RevParseVerified(dir, ref)
	if err != nil {
		return false, err
	}
	headCommit, err := revParseVerified(dir, head)
	if err != nil {
		return false, err
	}
	if headCommit == refCommit {
		return true, nil
	}

	refParentCommit, err := RevParseVerified(dir, refCommit+"^1")
	if err != nil {
		return false, err
	}

	output, err := cmd.Dir(dir).Output("git", "rev-list", "--first-parent", headCommit, "^"+refParentCommit)
	if err != nil {
		return false, err
	}
	commits := splitLinesFilterEmpty(output)
	for _, commit := range commits {
		if commit != refCommit {
			continue
		}
		return true, nil
	}

	return false, nil
}
