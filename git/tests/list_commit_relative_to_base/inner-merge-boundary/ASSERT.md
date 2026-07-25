## Expected

- `hidden`, the inner merge's P2, is absent.

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
        t.Fatalf("hashes=%v, want non-recursive recovery %v", resp.Hashes, req.Expected)
    }
}
```
