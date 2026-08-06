## Steps
- Set RefA to `main` (a valid ref)
- Set RefB to a non-existent reference `nonexistent-branch`

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.RefA = "main"
	req.RefB = "nonexistent-branch"
	return nil
}
```
