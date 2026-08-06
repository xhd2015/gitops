# IsInsideGit Test Case Tree

Run with:
```sh
doctest test ./ -v
```

## Test Case Index

| # | Path | Preconditions | Expected |
|---|---|---|---|
| 1 | `inside-git/` | `dir` is a git repo (initialised) | `(true, nil)` |
| 2 | `outside-git/` | `dir` is NOT a git repo | `(false, nil)` |

## Root Run

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

type Request struct {
    Dir string
}

type Response struct {
    Inside bool
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
    inside, err := IsInsideGit(req.Dir)
    if err != nil {
        return nil, err
    }
    return &Response{Inside: inside}, nil
}
```
