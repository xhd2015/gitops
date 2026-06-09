## Expected
- Result contains `.gitignore` (untracked file) and `_base/base.go` (from parent has-dirs setup)
- `build/` is not reported because all contents match `*.o` in `.gitignore`

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
	if len(resp.Files) != 2 {
		t.Fatalf("expected exactly 2 files, got %d: %v", len(resp.Files), resp.Files)
	}
}
```
