# Scenario

**Feature**: a gitignored directory is listed by ListIgnoredDirs

## Preconditions
- Root `.gitignore` contains `ignored/`
- `ignored/file.txt` exists in the repo

## Steps
1. Run `git init` and `git branch -M master` in dir.
2. Write `.gitignore`=`ignored/` and `ignored/file.txt`.
3. Stage and commit (the ignored dir stays untracked).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir
	exec.Command("git", "-C", dir, "init").Run()
	exec.Command("git", "-C", dir, "branch", "-M", "master").Run()
	os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("ignored/\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "ignored"), 0755)
	os.WriteFile(filepath.Join(dir, "ignored", "file.txt"), []byte("x"), 0644)
	exec.Command("git", "-C", dir, "add", ".").Run()
	exec.Command("git", "-C", dir, "commit", "-m", "init").Run()
	return nil
}
```
