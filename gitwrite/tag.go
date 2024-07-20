package gitwrite

import (
	"fmt"

	"github.com/xhd2015/xgo/support/cmd"
)

func Tag(dir string, tag string) error {
	if tag == "" {
		return fmt.Errorf("requires tag")
	}
	return cmd.Dir(dir).Run("git", "tag", tag)
}
