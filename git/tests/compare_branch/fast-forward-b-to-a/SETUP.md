## Steps
- Same setup as a-to-b but with RefA and RefB swapped

```go
import (
	"os"
	"testing"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir
	runGit := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}

	runGit("checkout", "-b", "a")
	runGit("checkout", "main")
	err := os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("content 1"), 0644)
	if err != nil {
		return err
	}
	runGit("add", ".")
	runGit("commit", "-m", "commit on main")

	req.RefA = "main"
	req.RefB = "a"
	return nil
}
```
