## Expected

- `FetchTargetedRefs` succeeds.
- Dest ref is `refs/gitops/targets/<64 hex>` — a single slash between `targets` and the hash.
- Dest is **not** `refs/gitops/targets//<hash>`.

## Side Effects

- One dest ref under the normalized default prefix.

## Errors

- None.

## Exit Code

- Exit code is zero.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
	if err != nil {
		t.Fatalf("FetchTargetedRefs: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	hash := destHashOf("master")
	want := DefaultTargetRefPrefix + "/" + hash
	doubled := DefaultTargetRefPrefix + "//" + hash
	if !hasRef(resp.AllRefs, want) {
		t.Fatalf("missing normalized dest %s; refs=%v", want, resp.AllRefs)
	}
	if hasRef(resp.AllRefs, doubled) {
		t.Fatalf("trailing-slash prefix produced double-slash dest %s; refs=%v", doubled, resp.AllRefs)
	}
	for _, r := range resp.AllRefs {
		if strings.Contains(r, "//") {
			t.Fatalf("dest ref has double slash: %s", r)
		}
	}
}
```
