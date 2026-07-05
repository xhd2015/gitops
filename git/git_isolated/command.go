package git_isolated

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Command builds a git subprocess with isolated env and default test identity.
func Command(dir string, args ...string) *exec.Cmd {
	return DefaultConfig.Command(dir, args...)
}

// CommandWith builds a git subprocess using cfg.
func (c Config) Command(dir string, args ...string) *exec.Cmd {
	allArgs := append([]string{
		"-c", "user.email=" + c.userEmail(),
		"-c", "user.name=" + c.userName(),
	}, args...)
	cmd := exec.Command("git", allArgs...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), Env()...)
	if len(c.ExtraEnv) > 0 {
		cmd.Env = append(cmd.Env, c.ExtraEnv...)
	}
	return cmd
}

// Run runs git in dir and returns an error on non-zero exit.
func Run(dir string, args ...string) error {
	return DefaultConfig.Run(dir, args...)
}

// Output runs git in dir and returns trimmed stdout.
func Output(dir string, args ...string) (string, error) {
	return DefaultConfig.Output(dir, args...)
}

// CombinedOutput runs git in dir and returns combined stdout/stderr.
func CombinedOutput(dir string, args ...string) ([]byte, error) {
	return DefaultConfig.CombinedOutput(dir, args...)
}

// Run runs git in dir and returns an error on non-zero exit.
func (c Config) Run(dir string, args ...string) error {
	_, err := c.CombinedOutput(dir, args...)
	return err
}

// Output runs git in dir and returns trimmed stdout.
func (c Config) Output(dir string, args ...string) (string, error) {
	out, err := c.CombinedOutput(dir, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// CombinedOutput runs git in dir and returns combined stdout/stderr.
func (c Config) CombinedOutput(dir string, args ...string) ([]byte, error) {
	out, err := c.Command(dir, args...).CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("git %v: %w\n%s", args, err, out)
	}
	return out, nil
}