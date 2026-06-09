## Expected
- Result contains `_utnew.go` (from parent untracked setup) and `view/a.go`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	found := make(map[string]bool)
	for _, f := range resp.Files {
		found[f] = true
	}
	if !found["_utnew.go"] {
		t.Fatalf("expected _utnew.go in result, got: %v", resp.Files)
	}
	if !found["view/a.go"] {
		t.Fatalf("expected view/a.go in result, got: %v", resp.Files)
	}
	if len(resp.Files) != 2 {
		t.Fatalf("expected exactly 2 files, got %d: %v", len(resp.Files), resp.Files)
	}
}
```
