## Expected
- An error is returned because the directory is not a git repository

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for non-git directory, got nil")
	}
}
```
