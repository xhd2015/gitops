## Expected

- `git_isolated.Init` + add + commit succeeds.
- Latest commit author email is `test@test.com`.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/gitops/git/git_isolated"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cfg := git_isolated.Config{ExtraEnv: []string{"GIT_CONFIG_GLOBAL=" + req.GlobalConfig}}
	if err := cfg.Run(req.RepoDir, "-c", "init.templateDir=", "init", "-b", "main"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(req.RepoDir, "README.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Command(req.RepoDir, "add", "README.md").Run(); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Command(req.RepoDir, "commit", "-m", "init").Run(); err != nil {
		t.Fatal(err)
	}
	email, err := cfg.Output(req.RepoDir, "log", "-1", "--pretty=%ae")
	if err != nil {
		t.Fatal(err)
	}
	if email != git_isolated.DefaultUserEmail {
		t.Fatalf("email = %q, want %q", email, git_isolated.DefaultUserEmail)
	}
}
```