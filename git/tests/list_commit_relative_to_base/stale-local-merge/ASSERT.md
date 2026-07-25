## Expected

- `698` and its basic P1 commit `19b` are present.
- Its direct P2 first-parent chain `2b -> 61 -> 9f` is present exactly once.

```go
import (
    "reflect"
    "testing"

    "github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatal(err)
    }
    if !reflect.DeepEqual(resp.Hashes, req.Expected) {
        t.Fatalf("hashes=%v, want recovered P2 chain %v", resp.Hashes, req.Expected)
    }
}
```
