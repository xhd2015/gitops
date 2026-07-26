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
    if !reflect.DeepEqual(resp.RelativeHashes, req.Expected) {
        t.Fatalf("relative hashes=%v, want recovered P2 chain %v", resp.RelativeHashes, req.Expected)
    }
    if !reflect.DeepEqual(resp.DiffPointHashes, req.Expected) {
        t.Fatalf("diff-point hashes=%v, want recovered P2 chain %v", resp.DiffPointHashes, req.Expected)
    }
}
```
