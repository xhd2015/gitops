## Expected
- Result is nil (deleted files are filtered out)
- ResolvePathsToFiles is not set, result reflects raw git status output

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 0 {
		t.Fatalf("expected nil (deleted excluded), got: %v", resp.Files)
	}
}
```
