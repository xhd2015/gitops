## Expected

- `GetStagedFiles` returns `["sub/pkg.go"]` — the path is relative to the repo root.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 {
		t.Fatalf("expected 1 staged file, got %d: %v", len(resp.Files), resp.Files)
	}
	if resp.Files[0] != "sub/pkg.go" {
		t.Fatalf("expected sub/pkg.go, got %q", resp.Files[0])
	}
}
```
