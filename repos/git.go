package repos

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func runGit(dir string, auth *GitAuthConfig, args ...string) error {
	return runGitWithEnv(dir, auth, nil, args...)
}

// runGitWithEnv runs git and, on failure, attaches CombinedOutput (stdout+stderr)
// so callers surface fatal/invalid-ref diagnostics instead of bare exit status.
func runGitWithEnv(dir string, auth *GitAuthConfig, env map[string]string, args ...string) error {
	return runGitWithEnvContext(context.Background(), dir, auth, env, args...)
}

func runGitWithEnvContext(ctx context.Context, dir string, auth *GitAuthConfig, env map[string]string, args ...string) error {
	_, err := gitOutputWithEnvContext(ctx, dir, auth, env, args...)
	return err
}

func gitOutputWithEnvContext(ctx context.Context, dir string, auth *GitAuthConfig, env map[string]string, args ...string) (string, error) {
	out, err := gitOutputBytesWithEnvContext(ctx, dir, auth, env, args...)
	return strings.TrimSpace(string(out)), err
}

func gitOutputBytesWithEnvContext(ctx context.Context, dir string, auth *GitAuthConfig, env map[string]string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "git", auth.gitArgs(args...)...)
	command.Dir = dir
	if len(env) > 0 {
		command.Env = MergeGitEnv(env, nil)
	}
	out, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s in %s: %w\n%s", strings.Join(args, " "), dir, err, out)
	}
	return out, nil
}

func gitStdoutWithEnvContext(ctx context.Context, dir string, auth *GitAuthConfig, env map[string]string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "git", auth.gitArgs(args...)...)
	command.Dir = dir
	if len(env) > 0 {
		command.Env = MergeGitEnv(env, nil)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("git %s in %s: %w\n%s", strings.Join(args, " "), dir, err, stderr.Bytes())
	}
	return stdout.Bytes(), nil
}

// runGitStreamStderr runs git with stdout+stderr attached to w so progress
// (e.g. fetch --progress) appears live. Does not use CombinedOutput.
func runGitStreamStderr(ctx context.Context, dir string, auth *GitAuthConfig, env map[string]string, w io.Writer, args ...string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	command := exec.CommandContext(ctx, "git", auth.gitArgs(args...)...)
	command.Dir = dir
	if len(env) > 0 {
		command.Env = MergeGitEnv(env, nil)
	}
	if w == nil {
		w = io.Discard
	}
	command.Stdout = w
	command.Stderr = w
	if err := command.Run(); err != nil {
		return fmt.Errorf("git %s in %s: %w", strings.Join(args, " "), dir, err)
	}
	return nil
}

func gitOutput(dir string, args ...string) (string, error) {
	command := exec.Command("git", args...)
	command.Dir = dir
	out, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s in %s: %w\n%s", strings.Join(args, " "), dir, err, out)
	}
	return strings.TrimSpace(string(out)), nil
}
