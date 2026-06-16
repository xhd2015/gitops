## Preconditions
- CompareWith option is set to a valid commit ref by child tests

## Steps
1. Ensure ResolvePaths defaults to false; child tests will override as needed

```go
func Setup(t *testing.T, req *Request) error {
	req.ResolvePaths = false
	return nil
}
```
