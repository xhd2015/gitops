## Expected

- `GetStagedFiles` returns an empty slice.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 0 {
		t.Fatalf("expected 0 staged files, got %d: %v", len(resp.Files), resp.Files)
	}
}
```
