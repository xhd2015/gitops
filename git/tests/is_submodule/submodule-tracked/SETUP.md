# Scenario

**Feature**: IsSubmodule returns true for a path tracked as a real submodule gitlink

## Preconditions
- `ext/` is tracked by `dir` as a real submodule (gitlink mode 160000).

## Steps
1. Create a separate child git repo (`child`).
2. In `dir`, run `git submodule add <child> ext` so `ext/` becomes a tracked
   submodule gitlink, then commit.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir

	// Build a separate child repo to wire in as a submodule.
	child, err := os.MkdirTemp("", "gitchild")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(child) })
	exec.Command("git", "-C", child, "init").Run()
	exec.Command("git", "-C", child, "branch", "-M", "master").Run()
	exec.Command("git", "-C", child, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", child, "config", "user.name", "Test User").Run()
	os.WriteFile(filepath.Join(child, "go.mod"), []byte("module example.com/ext\n\ngo 1.22\n"), 0644)
	exec.Command("git", "-C", child, "add", ".").Run()
	exec.Command("git", "-C", child, "commit", "-m", "init").Run()

	if out, err := exec.Command("git", "-C", dir, "-c", "protocol.file.allow=always", "submodule", "add", child, "ext").CombinedOutput(); err != nil {
		t.Fatalf("git submodule add: %v\n%s", err, out)
	}
	exec.Command("git", "-C", dir, "commit", "-m", "add submodule").Run()
	return nil
}
```
