## Expected
- Result is nil (no file changes between HEAD~1 and HEAD, both clean commits)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 0 {
		t.Fatalf("expected nil (no diff between clean commits), got: %v", resp.Files)
	}
}
```
