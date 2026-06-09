## Steps
1. Set req.Dir to a path that does not exist on disk

```go
func Setup(t *testing.T, req *Request) error {
	req.Dir = "/nonexistent/path/that/does/not/exist"
	return nil
}
```
