## Preconditions
- ResolvePathsToFiles option is enabled

## Steps
1. Set req.ResolvePaths = true

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.ResolvePaths = true
	return nil
}
```
