## Expected
- Result is `[_base/base.go]` (from parent has-dirs setup). Git does not report the `build/` dir because all contents are gitignored.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "_base/base.go" {
		t.Fatalf("expected [_base/base.go], got: %v", resp.Files)
	}
}
```
