## Expected
- `CheckIgnore` returns `false` (`important.o` is un-ignored by `!important.o`)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.Ignored {
		t.Fatal("expected important.o to NOT be gitignored (negation), got true")
	}
}
```
