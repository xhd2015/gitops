package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// RunGit runs `git <args...>` with cmd.Dir=dir and returns trimmed stdout.
// Arguments are passed as argv — never through a shell — so values like
// "HEAD|id" cannot trigger shell metacharacter evaluation.
func RunGit(dir string, args ...string) (string, error) {
	return runGit(dir, nil, args...)
}

// RunGitAllowExit is like RunGit but treats the listed process exit codes as
// success (stdout still returned). Useful for `git grep`, which exits 1 when
// there are no matches.
func RunGitAllowExit(dir string, okExits []int, args ...string) (string, error) {
	return runGit(dir, okExits, args...)
}

func runGit(dir string, okExits []int, args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("requires git args")
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	out := strings.TrimSuffix(stdout.String(), "\n")
	if err == nil {
		return out, nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		code := exitErr.ExitCode()
		for _, c := range okExits {
			if code == c {
				return out, nil
			}
		}
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return out, fmt.Errorf("git %v: exit %d: %s", args, code, msg)
		}
		return out, fmt.Errorf("git %v: exit %d", args, code)
	}
	return out, fmt.Errorf("git %v: %w", args, err)
}
