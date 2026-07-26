## Expected

- Every direct side-chain commit is returned once, even where paths meet at `a`.

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
    if !reflect.DeepEqual(resp.RelativeHashes, req.Expected) || !reflect.DeepEqual(resp.DiffPointHashes, req.Expected) {
        t.Fatalf("relative=%v diff-points=%v, want deduplicated recovered chains %v", resp.RelativeHashes, resp.DiffPointHashes, req.Expected)
    }
}
```
