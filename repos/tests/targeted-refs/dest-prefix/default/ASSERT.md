## Expected

- `FetchTargetedRefs` succeeds.
- Dest ref is exactly `refs/gitops/targets/` + 64-char lowercase hex(sha256(`master`)).
- No dest exists under `refs/other/targets/`.
- `Commits["master"]` is a 40-char hex SHA.
- Cache dir is under the injected `ReposRoot` and includes `targeted-cache`.

## Side Effects

- Partial bare created under `{ReposRoot}/targeted-cache/…`.
- Dest ref created only under the default gitops prefix.

## Errors

- None.

## Exit Code

- Exit code is zero.

```go
import (
	"path/filepath"
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
	want := expectedDestRef("", "master")
	if !strings.HasPrefix(want, DefaultTargetRefPrefix+"/") {
		t.Fatalf("test helper dest %q not under default prefix", want)
	}
	if !hasRef(resp.AllRefs, want) {
		t.Fatalf("missing dest %s; refs=%v", want, resp.AllRefs)
	}
	if hasRefPrefix(resp.AllRefs, CustomTargetRefPrefix+"/") {
		t.Fatalf("default prefix must not write %s/…; refs=%v", CustomTargetRefPrefix, resp.AllRefs)
	}
	if resp.Commits["master"] == "" || len(resp.Commits["master"]) != 40 {
		t.Fatalf("Commits[master]=%q, want 40-hex", resp.Commits["master"])
	}
	if req.ReposRoot == "" || resp.CacheDir == "" {
		t.Fatal("expected injected ReposRoot and CacheDir")
	}
	if !strings.HasPrefix(filepath.Clean(resp.CacheDir), filepath.Clean(req.ReposRoot)) {
		t.Fatalf("CacheDir %q not under ReposRoot %q", resp.CacheDir, req.ReposRoot)
	}
	if !strings.Contains(filepath.ToSlash(resp.CacheDir), "/targeted-cache/") {
		t.Fatalf("CacheDir %q missing targeted-cache segment", resp.CacheDir)
	}
}
```
