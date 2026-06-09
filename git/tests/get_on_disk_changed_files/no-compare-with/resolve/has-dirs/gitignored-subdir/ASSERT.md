## Expected
- Result contains only `view/a.go`
- `view/build/output.js` is excluded because `build/` is in `.gitignore`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	found := make(map[string]bool)
	for _, f := range resp.Files {
		found[f] = true
	}
	if !found["view/a.go"] {
		t.Fatalf("expected view/a.go in result, got: %v", resp.Files)
	}
	if found["view/build/output.js"] {
		t.Fatalf("view/build/output.js should be excluded (build/ in .gitignore), got: %v", resp.Files)
	}
	if len(resp.Files) != 1 {
		t.Fatalf("expected exactly 1 file, got %d: %v", len(resp.Files), resp.Files)
	}
}
```
