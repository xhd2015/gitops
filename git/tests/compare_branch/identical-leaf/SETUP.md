## Steps
- Create a tag `v1` pointing to the same commit as `main`
- Set RefA to `main` and RefB to `v1` — both resolve to the same commit

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"os/exec"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	cmd := exec.Command("git", "tag", "v1")
	cmd.Dir = req.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git tag failed: %v\n%s", err, out)
	}

	req.RefA = "main"
	req.RefB = "v1"
	return nil
}
```
