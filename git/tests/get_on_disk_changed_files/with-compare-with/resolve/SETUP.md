## Preconditions
- ResolvePathsToFiles option is enabled

## Steps
1. Set req.ResolvePaths = true

```go
func Setup(t *testing.T, req *Request) error {
	req.ResolvePaths = true
	return nil
}
```
