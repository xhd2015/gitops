## Expected
- Result contains `mod.go` and `new.go` (del.go excluded, keep.go unchanged)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	foundMod, foundNew := false, false
	for _, f := range resp.Files {
		if f == "mod.go" {
			foundMod = true
		}
		if f == "new.go" {
			foundNew = true
		}
	}
	if !foundMod || !foundNew || len(resp.Files) != 2 {
		t.Fatalf("expected [mod.go, new.go], got: %v", resp.Files)
	}
}
```
