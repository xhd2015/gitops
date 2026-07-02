## Expected

- `ShowToplevel(<repo>/go-pkgs)` returns `<repo>\n`.
- The output keeps the trailing newline from Git's raw command-style output.
- No error is returned.

## Side Effects

- No repository files are modified.

## Errors

- None.

## Exit Code

- The doctest assertion succeeds.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TopOutput != req.ExpectedTop {
		t.Fatalf("ShowToplevel without GIT_DIR = %q, want %q", resp.TopOutput, req.ExpectedTop)
	}
}
```
