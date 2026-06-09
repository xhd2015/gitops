## Expected
- An error is returned because git commands fail against the corrupted repository

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for corrupted git repo, got nil")
	}
}
```
