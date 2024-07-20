package git

import (
	"fmt"

	"github.com/xhd2015/xgo/support/cmd"
)

func ListFile(dir string, ref string) ([]string, error) {
	if ref == "" {
		return nil, fmt.Errorf("requires revision")
	}
	args := []string{"ls-files"}
	if ref != COMMIT_WORKING {
		args = append(args, "--with-tree", ref)
	}
	res, err := cmd.Dir(dir).Output("git", args...)
	if err != nil {
		return nil, err
	}
	return splitLinesFilterEmpty(res), nil
}
