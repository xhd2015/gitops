## Expected
- Result contains `new.go` (new rename target)
- ResolvePathsToFiles is not set, result reflects raw git status output

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 1 || resp.Files[0] != "new.go" {
		t.Fatalf("expected [new.go], got: %v", resp.Files)
	}
}
```
