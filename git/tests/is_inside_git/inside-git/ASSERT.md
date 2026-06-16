## Expected
- `resp.Inside` is `true`
- No error

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !resp.Inside {
        t.Fatal("expected Inside=true, got false")
    }
}
```
