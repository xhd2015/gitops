## Preconditions
- ResolvePathsToFiles option is NOT enabled

## Steps
1. Set req.ResolvePaths = false

```go
func Setup(t *testing.T, req *Request) error {
	req.ResolvePaths = false
	return nil
}
```
