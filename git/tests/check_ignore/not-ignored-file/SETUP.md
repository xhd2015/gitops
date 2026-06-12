## Preconditions
- Root `.gitignore` contains `*.log`
- `main.go` exists in the repo (does NOT match pattern)

## Steps
1. Copy `testdata/` fixtures into the repo directory
2. Stage and commit
3. Set `req.Path = "main.go"`

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	testdata := filepath.Join("testdata")
	entries, err := os.ReadDir(testdata)
	if err != nil {
		return err
	}
	for _, e := range entries {
		src := filepath.Join(testdata, e.Name())
		dst := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		os.WriteFile(dst, data, 0644)
	}
	exec.Command("git", "-C", dir, "add", ".").Run()
	exec.Command("git", "-C", dir, "commit", "-m", "add fixtures").Run()
	req.Path = "main.go"
	return nil
}
```
