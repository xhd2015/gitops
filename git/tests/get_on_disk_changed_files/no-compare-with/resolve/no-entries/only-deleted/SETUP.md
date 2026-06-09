## Steps
1. Run `git rm README.md` to stage a delete

```go
import "os/exec"

func Setup(t *testing.T, req *Request) error {
	return exec.Command("git", "-C", req.Dir, "rm", "README.md").Run()
}
```