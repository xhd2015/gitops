package gitwrite

import "github.com/xhd2015/xgo/support/cmd"

func AddAll(dir string) error {
	return cmd.Dir(dir).Run("git", "add", "-A")
}
