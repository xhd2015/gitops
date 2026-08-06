## Expected

- `RestoreStaged` succeeds (no error).
- Neither `a.txt` nor `b.txt` is in the staging area.
- Both files still exist on disk.

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"path/filepath"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range resp.StagedAfter {
		if f == "a.txt" || f == "b.txt" {
			t.Fatalf("expected %s to be unstaged, but it is still in the index", f)
		}
	}
	for _, name := range []string{"a.txt", "b.txt"} {
		if !fileExists(filepath.Join(req.Dir, name)) {
			t.Fatalf("expected %s to still exist on disk after unstage, but it does not", name)
		}
	}
}
```
