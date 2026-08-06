## Preconditions
- ResolvePathsToFiles option is NOT enabled

## Steps
1. Set req.ResolvePaths = false

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.ResolvePaths = false
	return nil
}
```
