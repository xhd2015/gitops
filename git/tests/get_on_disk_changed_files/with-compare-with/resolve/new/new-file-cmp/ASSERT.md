## Expected
- Result contains `staged.go`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "staged.go" {
		t.Fatalf("expected [staged.go], got: %v", resp.Files)
	}
}
```
