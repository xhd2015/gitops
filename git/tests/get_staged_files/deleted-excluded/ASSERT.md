## Expected

- `GetStagedFiles` returns an empty slice (deleted files excluded by `--diff-filter=ACMRT`).

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 0 {
		t.Fatalf("expected 0 staged files when deleted file is staged, got %d: %v", len(resp.Files), resp.Files)
	}
}
```
