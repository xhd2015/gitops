## Errors

- An error is returned because the directory is not a git repository.
- Error need not be `*DiffCachedParseError` (ordinary git / exec failure is fine).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for non-git directory, got nil")
	}
}
```
