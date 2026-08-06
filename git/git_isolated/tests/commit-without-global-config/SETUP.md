# Scenario

Global gitconfig defines user.email but has no effect on isolated commits.

## Steps

1. Write a temp global gitconfig without usable identity for subprocess.
2. Init repo and commit via git_isolated with ExtraEnv pointing at that file.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	globalCfg := filepath.Join(t.TempDir(), "global-gitconfig")
	if err := os.WriteFile(globalCfg, []byte("[user]\n\temail = global@example.com\n"), 0o644); err != nil {
		return err
	}
	repo := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		return err
	}
	req.RepoDir = repo
	req.GlobalConfig = globalCfg
	return nil
}
```