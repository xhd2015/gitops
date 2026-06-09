## Steps
1. Run `git mv README.md new.go`

```go
import "os/exec"

func Setup(t *testing.T, req *Request) error {
	return exec.Command("git", "-C", req.Dir, "mv", "README.md", "new.go").Run()
}
```