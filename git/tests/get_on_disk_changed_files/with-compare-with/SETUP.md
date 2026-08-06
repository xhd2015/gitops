## Preconditions
- CompareWith option is set to a valid commit ref by child tests

## Steps
1. Ensure ResolvePaths defaults to false; child tests will override as needed

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
