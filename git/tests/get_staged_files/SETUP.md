# Scenario

**Feature**: GetStagedFiles returns staged file paths from the git index

```
# caller passes a repo dir, GetStagedFiles runs git diff --cached
caller -> GetStagedFiles(dir) -> git diff --cached --name-only --diff-filter=ACMRT -> file list

# deleted files are excluded by --diff-filter=ACMRT
git index (ACMRT filter) -> excludes D status entries -> staged file paths
```

## Preconditions

- `git` is available in PATH
- `dir` is a git repository with an initial commit (required for `git rm` to work)

## Steps

1. Create a temporary directory and git-init it with an initial commit.
2. Each leaf creates and stages specific files.
3. Run calls `GetStagedFiles(req.Dir)` and returns the file list.

## Context

- Go module: `github.com/xhd2015/gitops`
- Package under test: `git`
- `GetStagedFiles(dir)` returns `([]string, error)` — staged file paths relative to `dir`.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir, err := os.MkdirTemp("", "gittest")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	exec.Command("git", "-C", dir, "init").Run()
	exec.Command("git", "-C", dir, "branch", "-M", "master").Run()
	exec.Command("git", "-C", dir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", dir, "config", "user.name", "Test User").Run()
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("test"), 0644)
	exec.Command("git", "-C", dir, "add", ".").Run()
	exec.Command("git", "-C", dir, "commit", "-m", "init").Run()
	req.Dir = dir
	return nil
}

func writeFile(dir string, name string, content string) error {
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func stageFile(dir string, name string, content string) error {
	if err := writeFile(dir, name, content); err != nil {
		return err
	}
	return exec.Command("git", "-C", dir, "add", name).Run()
}

func stagedFileNames(t *testing.T, dir string) []string {
	out, err := exec.Command("git", "-C", dir, "diff", "--cached", "--name-only", "--diff-filter=ACMRT").Output()
	if err != nil {
		t.Fatalf("git diff --cached --name-only: %v", err)
	}
	var files []string
	for _, line := range splitLines(string(out)) {
		if line != "" {
			files = append(files, line)
		}
	}
	return files
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
```
