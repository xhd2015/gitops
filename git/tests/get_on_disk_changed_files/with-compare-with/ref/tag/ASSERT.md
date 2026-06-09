## Expected
- Result contains `base.go` and `a.go` (base.go from parent ref setup, a.go added after tag)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	found := make(map[string]bool)
	for _, f := range resp.Files {
		found[f] = true
	}
	if !found["base.go"] {
		t.Fatalf("expected base.go in result, got: %v", resp.Files)
	}
	if !found["a.go"] {
		t.Fatalf("expected a.go in result, got: %v", resp.Files)
	}
	if len(resp.Files) != 2 {
		t.Fatalf("expected 2 files, got %d: %v", len(resp.Files), resp.Files)
	}
}
```
