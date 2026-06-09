## Expected
- Result is nil (clean tree vs same branch = no diff)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 0 {
		t.Fatalf("expected nil, got: %v", resp.Files)
	}
}
```
