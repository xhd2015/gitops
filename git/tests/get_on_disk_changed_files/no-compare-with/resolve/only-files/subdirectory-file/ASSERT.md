## Expected
- Result contains `sub/pkg/foo.go`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "sub/pkg/foo.go" {
		t.Fatalf("expected [sub/pkg/foo.go], got: %v", resp.Files)
	}
}
```
