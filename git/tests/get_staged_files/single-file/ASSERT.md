## Expected

- `GetStagedFiles` returns `["file.txt"]`.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 {
		t.Fatalf("expected 1 staged file, got %d: %v", len(resp.Files), resp.Files)
	}
	if resp.Files[0] != "file.txt" {
		t.Fatalf("expected file.txt, got %q", resp.Files[0])
	}
}
```
