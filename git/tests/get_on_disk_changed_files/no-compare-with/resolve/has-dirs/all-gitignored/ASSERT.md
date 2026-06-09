## Expected
- Result is nil (empty). Git does not report the directory at all because all contents are gitignored.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 0 {
		t.Fatalf("expected nil/empty result (git ignores dir with only ignored files), got: %v", resp.Files)
	}
}
```
