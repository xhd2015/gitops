## Expected
- An error is returned because the directory does not exist

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for nonexistent directory, got nil")
	}
}
```
