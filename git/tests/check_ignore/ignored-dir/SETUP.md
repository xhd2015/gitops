## Preconditions
- Root `.gitignore` contains `build/`
- `build/output.js` exists in the repo

## Steps
1. Copy `testdata/` fixtures into the repo directory
2. Stage and commit
3. Set `req.Path = "build"`

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func copyTestdata(t *testing.T, dir, testdata string) {
	t.Helper()
	filepath.WalkDir(testdata, func(src string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(testdata, src)
		if rel == "." {
			return nil
		}
		dst := filepath.Join(dir, rel)
		if d.IsDir() {
			return os.MkdirAll(dst, 0755)
		}
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, 0644)
	})
}

func Setup(t *testing.T, req *Request) error {
	dir := req.Dir
	testdata := filepath.Join("testdata")
	copyTestdata(t, dir, testdata)
	exec.Command("git", "-C", dir, "add", ".").Run()
	exec.Command("git", "-C", dir, "commit", "-m", "add fixtures").Run()
	req.Path = "build"
	return nil
}
```
