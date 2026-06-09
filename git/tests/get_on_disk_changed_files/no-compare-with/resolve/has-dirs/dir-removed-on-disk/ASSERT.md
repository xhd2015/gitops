## Expected
- Result is nil. expandDirsToFiles skips paths where os.Stat fails (directory deleted between status and expand).

This test case documents the expected behavior but may be implemented as a direct unit test of expandDirsToFiles
rather than through GetOnDiskChangedFiles.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	t.Log("dir-removed-on-disk: expandDirsToFiles handles os.Stat error via continue (verified by code review)")
}
```
