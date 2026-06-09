## Preconditions
- CompareWith option is NOT set

## Steps
1. Call GetOnDiskChangedFiles(dir) without CompareWith (uses git status --porcelain internally)

```go
func Run(t *testing.T, req *Request) (*Response, error) {
	var opts []onDiskChangedFileOption
	if req.ResolvePaths {
		opts = append(opts, ResolvePathsToFiles())
	}
	files, err := GetOnDiskChangedFiles(req.Dir, opts...)
	if err != nil {
		return nil, err
	}
	return &Response{Files: files}, nil
}
```
