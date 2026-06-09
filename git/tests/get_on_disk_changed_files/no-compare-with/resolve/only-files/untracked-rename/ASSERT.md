## Expected
- Result contains `RENAMED.go` (old file deleted, new file untracked)
- `README.md` is NOT included (was removed)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "RENAMED.go" {
		t.Fatalf("expected [RENAMED.go], got: %v", resp.Files)
	}
}
```
