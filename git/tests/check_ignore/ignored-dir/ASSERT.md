## Expected
- `CheckIgnore` returns `true` (`build/` matches `build/` pattern in .gitignore)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Ignored {
		t.Fatal("expected build/ to be gitignored, got false")
	}
}
```
