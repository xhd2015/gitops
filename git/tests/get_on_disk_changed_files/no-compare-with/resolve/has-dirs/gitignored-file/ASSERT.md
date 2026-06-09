## Expected
- Result contains `.gitignore` (untracked file), `_base/base.go` (from parent has-dirs setup), and `view/a.go`
- `view/debug.log` is excluded because it matches `.gitignore` pattern `*.log`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	found := make(map[string]bool)
	for _, f := range resp.Files {
		found[f] = true
	}
	if !found[".gitignore"] {
		t.Fatalf("expected .gitignore in result, got: %v", resp.Files)
	}
	if !found["_base/base.go"] {
		t.Fatalf("expected _base/base.go in result, got: %v", resp.Files)
	}
	if !found["view/a.go"] {
		t.Fatalf("expected view/a.go in result, got: %v", resp.Files)
	}
	if found["view/debug.log"] {
		t.Fatalf("view/debug.log should be excluded (matches .gitignore), got: %v", resp.Files)
	}
	if len(resp.Files) != 3 {
		t.Fatalf("expected exactly 3 files, got %d: %v", len(resp.Files), resp.Files)
	}
}
```
