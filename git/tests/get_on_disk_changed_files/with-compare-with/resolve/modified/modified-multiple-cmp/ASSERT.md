## Expected
- Result contains `a.go` and `b.go`

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 2 {
		t.Fatalf("expected 2 files, got: %v", resp.Files)
	}
	foundA, foundB := false, false
	for _, f := range resp.Files {
		if f == "a.go" {
			foundA = true
		}
		if f == "b.go" {
			foundB = true
		}
	}
	if !foundA || !foundB {
		t.Fatalf("expected [a.go, b.go], got: %v", resp.Files)
	}
}
```
