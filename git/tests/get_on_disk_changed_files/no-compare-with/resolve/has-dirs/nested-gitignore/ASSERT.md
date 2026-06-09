## Expected
- Result contains `_base/base.go` (from parent has-dirs setup) and `view/sub/keep.go`
- `view/sub/a.tmp` is excluded because `view/sub/.gitignore` ignores `*.tmp`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	found := make(map[string]bool)
	for _, f := range resp.Files {
		found[f] = true
	}
	if !found["_base/base.go"] {
		t.Fatalf("expected _base/base.go in result, got: %v", resp.Files)
	}
	if !found["view/sub/keep.go"] {
		t.Fatalf("expected view/sub/keep.go in result, got: %v", resp.Files)
	}
	if found["view/sub/a.tmp"] {
		t.Fatalf("view/sub/a.tmp should be excluded (matches nested .gitignore *.tmp), got: %v", resp.Files)
	}
	if len(resp.Files) != 2 {
		t.Fatalf("expected exactly 2 files, got %d: %v", len(resp.Files), resp.Files)
	}
}
```
