## Preconditions
- ResolvePathsToFiles option is NOT enabled
- Directory paths are returned as-is (not expanded to files)

## Steps
1. Set req.ResolvePaths = false

```go
func Setup(t *testing.T, req *Request) error {
	req.ResolvePaths = false
	return nil
}
```
