## Expected

- `RestoreStaged` succeeds.
- `a.txt` is no longer in the staging area.
- `b.txt` is still in the staging area.
- Both files still exist on disk.

```go
import (
	"os"
	"path/filepath"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	hasA, hasB := false, false
	for _, f := range resp.StagedAfter {
		switch f {
		case "a.txt":
			hasA = true
		case "b.txt":
			hasB = true
		}
	}
	if hasA {
		t.Fatal("expected a.txt to be unstaged, but it is still in the index")
	}
	if !hasB {
		t.Fatal("expected b.txt to still be staged, but it was also unstaged")
	}
	for _, name := range []string{"a.txt", "b.txt"} {
		if !fileExists(filepath.Join(req.Dir, name)) {
			t.Fatalf("expected %s to still exist on disk after unstage, but it does not", name)
		}
	}
}
```
