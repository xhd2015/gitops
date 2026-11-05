package repos

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// Legacy helpers also have a finite lifetime. Explicit callers can supply a
// shorter deadline; no Git subprocess is allowed to outlive cancellation.
const defaultGitTimeout = 5 * time.Minute

// RunGitContext executes a local repository operation with bounded process-tree
// ownership. It does not add authentication or bypass repository access checks.
func RunGitContext(ctx context.Context, dir string, args ...string) error {
	return runGitWithEnvContext(ctx, dir, nil, nil, args...)
}

func runGit(dir string, auth *GitAuthConfig, args ...string) error {
	return runGitWithEnv(dir, auth, nil, args...)
}

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
	var output bytes.Buffer
	command := gitCommand(dir, auth, env, args...)
	command.Stdout = &output
	command.Stderr = &output
	if err := runGitCommand(ctx, command); err != nil {
		// auth arguments are never included in the printable command.
		return nil, fmt.Errorf("git %s in %s: %w\n%s", strings.Join(args, " "), dir, err, auth.maskError(fmt.Errorf("%s", output.Bytes())))
	}
	return output.Bytes(), nil
}

func gitStdoutWithEnvContext(ctx context.Context, dir string, auth *GitAuthConfig, env map[string]string, args ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	command := gitCommand(dir, auth, env, args...)
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := runGitCommand(ctx, command); err != nil {
		return nil, fmt.Errorf("git %s in %s: %w\n%s", strings.Join(args, " "), dir, err, auth.maskError(fmt.Errorf("%s", stderr.Bytes())))
	}
	return stdout.Bytes(), nil
}

func runGitStreamStderr(ctx context.Context, dir string, auth *GitAuthConfig, env map[string]string, w io.Writer, args ...string) error {
	if w == nil {
		w = io.Discard
	}
	command := gitCommand(dir, auth, env, args...)
	command.Stdout, command.Stderr = w, w
	if err := runGitCommand(ctx, command); err != nil {
		return fmt.Errorf("git %s in %s: %w", strings.Join(args, " "), dir, err)
	}
	return nil
}

func gitCommand(dir string, auth *GitAuthConfig, env map[string]string, args ...string) *exec.Cmd {
	command := exec.Command("git", auth.gitArgs(args...)...)
	command.Dir = dir
	if len(env) > 0 {
		command.Env = MergeGitEnv(env, nil)
	}
	return command
}

func runGitCommand(ctx context.Context, command *exec.Cmd) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, defaultGitTimeout)
	defer cancel()
	return runCommandTree(ctx, command)
}

func gitOutput(dir string, args ...string) (string, error) {
	return gitOutputWithEnvContext(context.Background(), dir, nil, nil, args...)
}
