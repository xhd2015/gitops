## Expected

- Two entries
- Paths `/path/to/repo` and `/path/to/repo-feature`
- Branches `master` and `feature` (refs/heads/ stripped)
- HEAD hashes preserved

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Entries) != 2 {
		t.Fatalf("len = %d, want 2", len(resp.Entries))
	}
	if resp.Entries[0].Path != "/path/to/repo" {
		t.Fatalf("entries[0].Path = %q", resp.Entries[0].Path)
	}
	if resp.Entries[0].Branch != "master" {
		t.Fatalf("entries[0].Branch = %q, want master", resp.Entries[0].Branch)
	}
	if resp.Entries[0].HEAD != "1234567890abcdef1234567890abcdef12345678" {
		t.Fatalf("entries[0].HEAD = %q", resp.Entries[0].HEAD)
	}
	if resp.Entries[1].Path != "/path/to/repo-feature" {
		t.Fatalf("entries[1].Path = %q", resp.Entries[1].Path)
	}
	if resp.Entries[1].Branch != "feature" {
		t.Fatalf("entries[1].Branch = %q, want feature", resp.Entries[1].Branch)
	}
	if resp.Entries[1].HEAD != "abcdef1234567890abcdef1234567890abcdef12" {
		t.Fatalf("entries[1].HEAD = %q", resp.Entries[1].HEAD)
	}
}
```
