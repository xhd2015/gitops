## Expected
- Result is `[a.go]` (only a.go added after tag v1.0; base.go was already committed before the tag)

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
