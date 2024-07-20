package gitwrite

import (
	"fmt"

	"github.com/xhd2015/gitops/git/worktree"

	"github.com/xhd2015/xgo/support/cmd"
)

func Commit(dir string, authorName string, authorMail string, msg string) error {
	if dir == "" {
		return fmt.Errorf("requires dir")
	}
	if authorName == "" {
		return fmt.Errorf("requires authorName")
	}
	if authorMail == "" {
		return fmt.Errorf("requires authorMail")
	}
	if msg == "" {
		return fmt.Errorf("requires msg")
	}
	clean, err := worktree.IndexClean(dir)
	if err != nil {
		return err
	}
	if clean {
		return fmt.Errorf("no changes made")
	}

	return cmd.Dir(dir).Run("git", "-c", "user.email="+authorMail, "-c", "user.name="+authorMail, "commit", "--author="+fmt.Sprintf("%s <%s>", authorName, authorMail), "-m", msg)
}
