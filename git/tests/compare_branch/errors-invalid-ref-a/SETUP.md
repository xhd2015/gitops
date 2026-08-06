## Steps
- Set RefA to a non-existent reference `nonexistent-branch`
- Set RefB to `main` (a valid ref)

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.RefA = "nonexistent-branch"
	req.RefB = "main"
	return nil
}
```
