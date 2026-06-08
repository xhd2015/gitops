package git

import (
	"strings"

	"github.com/xhd2015/xgo/support/cmd"
)

func GetOnDiskChangedFiles(dir string) ([]string, error) {
	output, err := cmd.Dir(dir).Output("git", "status", "--porcelain")
	if err != nil {
		return nil, err
	}

	rawLines := strings.Split(output, "\n")
	seen := make(map[string]bool, len(rawLines))
	var result []string
	for _, line := range rawLines {
		if len(line) < 4 {
			continue
		}
		xy := line[0:2]
		pathPart := line[3:]

		if xy[0] == 'D' || xy[1] == 'D' {
			continue
		}

		if xy[0] == 'R' {
			idx := strings.LastIndex(pathPart, " -> ")
			if idx >= 0 {
				pathPart = pathPart[idx+4:]
			}
		}

		pathPart = strings.TrimSpace(pathPart)
		if pathPart == "" {
			continue
		}
		if seen[pathPart] {
			continue
		}
		seen[pathPart] = true
		result = append(result, pathPart)
	}
	if result == nil {
		return nil, nil
	}
	return result, nil
}
