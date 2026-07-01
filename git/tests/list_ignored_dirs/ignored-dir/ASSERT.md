## Expected

- `ListIgnoredDirs` returns no error.
- The result contains exactly `ignored` (the ignored directory, slash-relative,
  no trailing slash, no `./`).

```go
import (
	"reflect"
	"sort"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	got := append([]string(nil), resp.Dirs...)
	sort.Strings(got)
	want := []string{"ignored"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListIgnoredDirs = %v, want %v", got, want)
	}
}
```
