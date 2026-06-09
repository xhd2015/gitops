## Expected
- Result contains `view/sub/keep.go`
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
	if !found["view/sub/keep.go"] {
		t.Fatalf("expected view/sub/keep.go in result, got: %v", resp.Files)
	}
	if found["view/sub/a.tmp"] {
		t.Fatalf("view/sub/a.tmp should be excluded (matches nested .gitignore *.tmp), got: %v", resp.Files)
	}
	if len(resp.Files) != 1 {
		t.Fatalf("expected exactly 1 file, got %d: %v", len(resp.Files), resp.Files)
	}
}
```
