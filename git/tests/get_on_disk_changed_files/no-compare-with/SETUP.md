## Preconditions
- CompareWith option is NOT set

## Steps
1. req.CompareWith is explicitly cleared — the root Run function will not pass CompareWith

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.CompareWith = ""
	return nil
}
```
