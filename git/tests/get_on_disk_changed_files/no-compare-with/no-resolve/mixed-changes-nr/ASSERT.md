## Expected
- Result contains `mod.go`, `added.go`, `untracked.go`, `ren_new.go`
- `del.go` is excluded (deleted) and `ren.go` is excluded (renamed away)
- ResolvePathsToFiles is not set, result reflects raw git status output

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	found := make(map[string]bool)
	for _, f := range resp.Files {
		found[f] = true
	}
	expected := []string{"mod.go", "added.go", "untracked.go", "ren_new.go"}
	for _, e := range expected {
		if !found[e] {
			t.Fatalf("expected %s in result, got: %v", e, resp.Files)
		}
	}
	if found["del.go"] {
		t.Fatalf("del.go should be excluded (deleted), got: %v", resp.Files)
	}
	if found["ren.go"] {
		t.Fatalf("ren.go should be excluded (renamed away), got: %v", resp.Files)
	}
	if len(resp.Files) != len(expected) {
		t.Fatalf("expected %d files, got %d: %v", len(expected), len(resp.Files), resp.Files)
	}
}
```
