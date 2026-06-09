## Expected
- Result is nil (deleted files are filtered out)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Files) != 0 {
		t.Fatalf("expected nil (deleted excluded), got: %v", resp.Files)
	}
}
```
