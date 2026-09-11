package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunGit_InjectionRefDoesNotExecuteShell(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)

	marker := filepath.Join(dir, "PWNED")
	// If this were passed through bash -c without proper quoting, `|touch…`
	// would create the marker file. With argv, git only sees a weird ref.
	payload := "HEAD|touch " + marker

	_, err := RunGit(dir, "rev-parse", payload+"^{commit}")
	if err == nil {
		t.Fatalf("expected invalid-ref error for payload %q", payload)
	}
	if _, statErr := os.Stat(marker); !os.IsNotExist(statErr) {
		t.Fatalf("shell metacharacters were evaluated; marker exists: %v", marker)
	}
}

func TestRunGitAllowExit_GrepNoMatch(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)
	mustWrite(t, filepath.Join(dir, "a.txt"), "hello\n")
	run(t, dir, "git", "add", "a.txt")
	run(t, dir, "git", "commit", "--no-verify", "-m", "add")

	out, err := RunGitAllowExit(dir, []int{1}, "grep", "-n", "-e", "nomatch-zzz", "HEAD", "--", "a.txt")
	if err != nil {
		t.Fatalf("exit 1 should be allowed: %v", err)
	}
	if strings.TrimSpace(out) != "" {
		t.Fatalf("expected empty stdout, got %q", out)
	}
}

func TestListCommits_UsesRefArgNotLiteral(t *testing.T) {
	dir := t.TempDir()
	initGitRepo(t, dir)
	mustWrite(t, filepath.Join(dir, "a.txt"), "v1\n")
	run(t, dir, "git", "add", "a.txt")
	run(t, dir, "git", "commit", "--no-verify", "-m", "c1")
	head, err := RunGit(dir, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}

	commits, err := ListCommits(dir, "", head)
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) == 0 {
		t.Fatalf("expected at least one commit for %s", head)
	}
	if commits[0].Hash == "" {
		t.Fatal("empty commit hash")
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	run(t, dir, "git", "init")
	run(t, dir, "git", "config", "user.email", "test@example.com")
	run(t, dir, "git", "config", "user.name", "test")
	run(t, dir, "git", "config", "core.hooksPath", "/dev/null")
	mustWrite(t, filepath.Join(dir, "README"), "x\n")
	run(t, dir, "git", "add", "README")
	run(t, dir, "git", "commit", "--no-verify", "-m", "init")
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func run(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out)
	}
}
