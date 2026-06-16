## Preconditions
- The function is expected to return an error

## Context
- Error cases are cross-cutting: they apply regardless of opts (CompareWith, ResolvePaths)

```go
func Setup(t *testing.T, req *Request) error {
	// Error cases: these tests expect the function to return an error.
	// Leave CompareWith and ResolvePaths at defaults — child tests
	// will override req.Dir or other fields as needed.
	req.CompareWith = ""
	req.ResolvePaths = false
	return nil
}
```
