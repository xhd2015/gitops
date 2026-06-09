## Expected
- Result contains `a/b/c/d.go`
- `a/b/ignored.log` is excluded because `*.log` is in root `.gitignore`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	found := make(map[string]bool)
	for _, f := range resp.Files {
		found[f] = true
	}
	if !found["a/b/c/d.go"] {
		t.Fatalf("expected a/b/c/d.go in result, got: %v", resp.Files)
	}
	if found["a/b/ignored.log"] {
		t.Fatalf("a/b/ignored.log should be excluded (matches *.log in .gitignore), got: %v", resp.Files)
	}
	if len(resp.Files) != 1 {
		t.Fatalf("expected exactly 1 file, got %d: %v", len(resp.Files), resp.Files)
	}
}
```
