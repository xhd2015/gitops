# Scenario

**Feature**: RestoreStaged removes files from the staging area without touching the working copy

```
# caller passes a repo dir and paths, runs git restore --staged
caller -> RestoreStaged(dir, paths...) -> git restore --staged -- <paths> -> files gone from index

# working copy files are preserved after unstage
git restore --staged -> removes from index only -> working copy unchanged
```

## Preconditions

- `git` is available in PATH
- `dir` is a git repository with an initial commit (required so `git restore --staged` can distinguish staged from committed content)

## Steps

1. Create a temporary directory and git-init it with an initial commit.
2. Each leaf creates and stages specific files, then sets `req.Paths`.
3. Run calls `RestoreStaged(req.Dir, req.Paths...)`, then returns the remaining staged files.

## Context

- Go module: `github.com/xhd2015/gitops`
- Package under test: `gitwrite`
- `RestoreStaged(dir, paths...)` runs `git restore --staged -- <paths>` in `dir`.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getStagedFileNames(dir string) ([]string, error) {
	cmd := exec.Command("git", "-C", dir, "diff", "--cached", "--name-only", "--diff-filter=ACMRT")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}
```
