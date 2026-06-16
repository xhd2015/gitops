## Preconditions
- `git` is available in PATH

## Steps
1. Each leaf sets `req.Dir` to the directory to test.
2. Run calls `IsInsideGit(req.Dir)` and returns the result.

## Context
- Go module: `github.com/xhd2015/gitops`
- Package under test: `git`
- `IsInsideGit(dir)` returns `(bool, error)` — true if dir is inside a git work tree.

```go
import (
    "testing"
)

type Request struct {
    Dir string
}

type Response struct {
    Inside bool
}

func Run(t *testing.T, req *Request) (*Response, error) {
    inside, err := IsInsideGit(req.Dir)
    if err != nil {
        return nil, err
    }
    return &Response{Inside: inside}, nil
}
```
