## Expected
- Result contains `a.go` (file added to HEAD relative to HEAD~1)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "a.go" {
		t.Fatalf("expected [a.go], got: %v", resp.Files)
	}
}
```
