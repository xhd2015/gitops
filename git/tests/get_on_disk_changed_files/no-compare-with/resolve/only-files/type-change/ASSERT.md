## Expected
- Result contains `file`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "file" {
		t.Fatalf("expected [file], got: %v", resp.Files)
	}
}
```
