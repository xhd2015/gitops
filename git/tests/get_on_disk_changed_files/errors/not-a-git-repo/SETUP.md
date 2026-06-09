## Steps
1. Set req.Dir to a directory that is NOT a git repository
2. Ensure the directory exists but has no .git subdirectory

```go
func Setup(t *testing.T, req *Request) error {
	dir, err := os.MkdirTemp("", "not-a-git-repo")
	if err != nil {
		return err
	}
	req.Dir = dir
	return nil
}
```
