## Expected

- Each `PorcelainCase` maps to the expected kind string.

```go
import (
	"testing"

	"github.com/xhd2015/gitops/git"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, tc := range req.PorcelainCases {
		got := git.TestExported_ClassifyPorcelainLine(tc.XY)
		if got != tc.Want {
			t.Fatalf("classifyPorcelainLine(%q) = %q, want %q", tc.XY, got, tc.Want)
		}
	}
}
```