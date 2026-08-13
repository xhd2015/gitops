## Expected

- `FetchTargetedRefs` succeeds.
- Dest ref is `refs/other/targets/` + the same 64-hex hash as the default prefix would use for input `master`.
- No dest exists under `refs/gitops/targets/`.

## Side Effects

- Dest written only under the injected custom prefix.

## Errors

- None.

## Exit Code

- Exit code is zero.

```go
import (
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
	want := expectedDestRef(CustomTargetRefPrefix, "master")
	if !hasRef(resp.AllRefs, want) {
		t.Fatalf("missing custom dest %s; refs=%v", want, resp.AllRefs)
	}
	defaultDest := expectedDestRef("", "master")
	if hasRef(resp.AllRefs, defaultDest) || hasRefPrefix(resp.AllRefs, DefaultTargetRefPrefix+"/") {
		t.Fatalf("custom prefix must not also write default dest; refs=%v", resp.AllRefs)
	}
}
```
