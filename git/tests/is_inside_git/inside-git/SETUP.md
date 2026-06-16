## Steps
1. Create a temp directory with `git init`
2. Set `req.Dir` to that directory

```go
import (
    "os"
    "os/exec"
)

func Setup(t *testing.T, req *Request) error {
    dir, err := os.MkdirTemp("", "isinsidegit")
    if err != nil {
        return err
    }
    t.Cleanup(func() { os.RemoveAll(dir) })
    if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
        t.Fatalf("git init: %v\n%s", err, out)
    }
    req.Dir = dir
    return nil
}
```
