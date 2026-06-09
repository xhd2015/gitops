## Expected
- Result contains `README.md` (deduplicated — both staged and unstaged appear once)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "README.md" {
		t.Fatalf("expected [README.md], got: %v", resp.Files)
	}
}
```
