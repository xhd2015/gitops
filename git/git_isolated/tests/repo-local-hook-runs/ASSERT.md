## Expected

- Commit fails because repo-local pre-commit hook runs.

```go
import (
	"os"
	"testing"

	"github.com/xhd2015/doctest/session"
	"path/filepath"
	"github.com/xhd2015/gitops/git/git_isolated"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := git_isolated.Init(req.RepoDir, "main"); err != nil {
		t.Fatal(err)
	}
	if err := git_isolated.Run(req.RepoDir, "config", "core.hooksPath", req.HooksDir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(req.RepoDir, "README.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := git_isolated.Run(req.RepoDir, "add", "README.md"); err != nil {
		t.Fatal(err)
	}
	if err := git_isolated.Run(req.RepoDir, "commit", "-m", "blocked"); err == nil {
		t.Fatal("expected repo-local pre-commit hook to block commit")
	}
}
```