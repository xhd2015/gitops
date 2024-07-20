package git

import (
	"github.com/xhd2015/gitops/git/worktree"
)

func WorkTreeClean(dir string) (ok bool, err error) {
	return worktree.IsClean(dir)
}

func IndexClean(dir string) (ok bool, err error) {
	return worktree.IndexClean(dir)
}
