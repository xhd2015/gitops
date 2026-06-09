## Expected
- Result contains `view/a.go` (DiffCommitFiles returns individual files, not dirs)
- CompareWith path produces file-level output, not directory entries

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "view/a.go" {
		t.Fatalf("expected [view/a.go], got: %v", resp.Files)
	}
}
```
