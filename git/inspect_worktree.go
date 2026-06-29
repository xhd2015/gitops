package git

import (
	"strings"

	"github.com/xhd2015/xgo/support/cmd"
)

// WorktreeInspect holds live git worktree state for a directory.
type WorktreeInspect struct {
	IsRepo        bool
	Branch        string // branch name, "(detached)", or "" when not a repo
	CommitShort   string // 7-char short hash
	CommitMessage string // subject line
	IsClean       bool
	Added         int
	Changed       int
	Renamed       int
	Deleted       int
	Uncommitted   int // total porcelain lines (for API backward compat)
}

// porcelainChangeType classifies a porcelain status line into a change bucket.
type porcelainChangeType int

const (
	porcelainNone porcelainChangeType = iota
	porcelainAdded
	porcelainChanged
	porcelainRenamed
	porcelainDeleted
)

// classifyPorcelainLine maps a porcelain XY status pair to a change type.
// Rules (first match wins):
//   - ?? → added
//   - R in either column → renamed
//   - D in either column → deleted
//   - A in either column → added
//   - M in either column → changed
func classifyPorcelainLine(xy string) porcelainChangeType {
	if len(xy) < 2 {
		return porcelainNone
	}
	if xy == "??" {
		return porcelainAdded
	}
	if xy[0] == 'R' || xy[1] == 'R' {
		return porcelainRenamed
	}
	if xy[0] == 'D' || xy[1] == 'D' {
		return porcelainDeleted
	}
	if xy[0] == 'A' || xy[1] == 'A' {
		return porcelainAdded
	}
	if xy[0] == 'M' || xy[1] == 'M' {
		return porcelainChanged
	}
	return porcelainNone
}

// InspectWorktree returns branch, HEAD commit, and porcelain change counts for dir.
// Non-git directories return a zero-valued result with IsRepo=false and no error.
func InspectWorktree(dir string) (*WorktreeInspect, error) {
	inside, err := IsInsideGit(dir)
	if err != nil {
		return nil, err
	}
	if !inside {
		return &WorktreeInspect{IsRepo: false, IsClean: true}, nil
	}

	result := &WorktreeInspect{IsRepo: true}

	branchOut, err := cmd.Dir(dir).Output("git", "branch", "--show-current")
	if err != nil {
		return nil, err
	}
	branch := strings.TrimSpace(branchOut)
	if branch == "" {
		result.Branch = "(detached)"
	} else {
		result.Branch = branch
	}

	commit, err := GetCommit(dir, "HEAD")
	if err == nil && commit != nil && commit.Hash != "" {
		hash := commit.Hash
		if len(hash) > 7 {
			hash = hash[:7]
		}
		result.CommitShort = hash
		result.CommitMessage = commit.Msg
	}

	statusOut, err := cmd.Dir(dir).Output("git", "status", "--porcelain=v1")
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(statusOut, "\n") {
		if len(line) < 3 {
			continue
		}
		xy := line[0:2]
		changeType := classifyPorcelainLine(xy)
		if changeType == porcelainNone {
			continue
		}
		result.Uncommitted++
		switch changeType {
		case porcelainAdded:
			result.Added++
		case porcelainChanged:
			result.Changed++
		case porcelainRenamed:
			result.Renamed++
		case porcelainDeleted:
			result.Deleted++
		}
	}

	result.IsClean = result.Uncommitted == 0
	return result, nil
}

// TestExported_ClassifyPorcelainLine exposes porcelain line classification for doctests.
func TestExported_ClassifyPorcelainLine(xy string) string {
	switch classifyPorcelainLine(xy) {
	case porcelainAdded:
		return "added"
	case porcelainChanged:
		return "changed"
	case porcelainRenamed:
		return "renamed"
	case porcelainDeleted:
		return "deleted"
	default:
		return "none"
	}
}