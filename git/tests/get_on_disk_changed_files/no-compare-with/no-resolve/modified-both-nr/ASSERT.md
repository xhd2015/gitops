## Expected
- Result contains `README.md` (deduplicated)
- ResolvePathsToFiles is not set, result reflects raw git status output

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
