## Steps
1. Set req.CompareWith = "nonexistent123"

```go
func Setup(t *testing.T, req *Request) error {
	req.CompareWith = "nonexistent123"
	return nil
}
```