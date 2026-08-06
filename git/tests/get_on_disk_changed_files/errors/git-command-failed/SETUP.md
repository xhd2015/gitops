## Steps
1. Create a git repository
2. Corrupt the .git directory so that git commands fail
3. For example, remove .git/HEAD or make .git unreadable

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Create a dir with a .git but no valid HEAD
	dir, err := os.MkdirTemp("", "broken-git")
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Join(dir, ".git"), 0755)
	req.Dir = dir
	return nil
}
```
