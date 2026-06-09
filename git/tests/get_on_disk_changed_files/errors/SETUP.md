## Preconditions
- The function is expected to return an error

## Context
- Error cases are cross-cutting: they apply regardless of opts (CompareWith, ResolvePaths)

```go
func Run(t *testing.T, req *Request) (*Response, error) {
	var opts []onDiskChangedFileOption
	if req.CompareWith != "" {
		opts = append(opts, CompareWith(req.CompareWith))
	}
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
