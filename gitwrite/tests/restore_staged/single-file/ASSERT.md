## Expected

- `RestoreStaged` succeeds (no error).
- `a.txt` is no longer in the staging area.
- `a.txt` still exists on disk.

## Side Effects

- `a.txt` is removed from the git index but preserved in the working tree.

```go
import (
	"os"
	"testing"

	"github.com/xhd2015/doctest/session"
)
import "path/filepath"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range resp.StagedAfter {
		if f == "a.txt" {
			t.Fatal("expected a.txt to be unstaged, but it is still in the index")
		}
	}
	aPath := filepath.Join(req.Dir, "a.txt")
	if !fileExists(aPath) {
		t.Fatal("expected a.txt to still exist on disk after unstage, but it does not")
	}
}
```
