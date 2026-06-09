## Expected
- Result is `[a.go]` (only a.go added after HEAD~1; base.go was already in HEAD~1)

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
