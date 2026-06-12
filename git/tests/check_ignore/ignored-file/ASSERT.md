## Expected
- `CheckIgnore` returns `true` (`app.log` matches `*.log` pattern)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Ignored {
		t.Fatal("expected app.log to be gitignored, got false")
	}
}
```
