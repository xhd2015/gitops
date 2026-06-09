## Expected
- Result contains `view/` (directory path, not expanded to files)
- ResolvePathsToFiles is not set, so directories are returned as-is

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "view/" {
		t.Fatalf("expected [view/] (dir not expanded), got: %v", resp.Files)
	}
}
```
