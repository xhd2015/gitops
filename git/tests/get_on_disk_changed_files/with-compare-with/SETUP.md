## Preconditions
- CompareWith option is set to a valid commit ref

## Steps
1. Call GetOnDiskChangedFiles(dir, CompareWith(ref))

```go
func Run(t *testing.T, req *Request) (*Response, error) {
	var opts []onDiskChangedFileOption
	opts = append(opts, CompareWith(req.CompareWith))
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
