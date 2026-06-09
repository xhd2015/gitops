## Preconditions
- ResolvePathsToFiles option is enabled

## Steps
1. Call GetOnDiskChangedFiles(dir, ResolvePathsToFiles())

```go
func Setup(t *testing.T, req *Request) error {
	req.ResolvePaths = true
	return nil
}
```
