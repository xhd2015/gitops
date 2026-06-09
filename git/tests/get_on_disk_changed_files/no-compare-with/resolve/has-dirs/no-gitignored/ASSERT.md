## Expected
- Result contains `_base/base.go` (from parent has-dirs setup), `view/a.go`, and `view/b.go`
- No other files included

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 3 {
		t.Fatalf("expected 3 files, got %d: %v", len(resp.Files), resp.Files)
	}
	found := make(map[string]bool)
	for _, f := range resp.Files {
		found[f] = true
	}
	if !found["_base/base.go"] {
		t.Fatalf("expected _base/base.go in result, got: %v", resp.Files)
	}
	if !found["view/a.go"] {
		t.Fatalf("expected view/a.go in result, got: %v", resp.Files)
	}
	if !found["view/b.go"] {
		t.Fatalf("expected view/b.go in result, got: %v", resp.Files)
	}
}
```
