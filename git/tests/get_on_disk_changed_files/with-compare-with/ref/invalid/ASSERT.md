## Expected
- Function returns an error for invalid ref

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected error for invalid ref, got nil")
	}
}
```
