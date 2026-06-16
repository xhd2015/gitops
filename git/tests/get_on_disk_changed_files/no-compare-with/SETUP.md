## Preconditions
- CompareWith option is NOT set

## Steps
1. req.CompareWith is explicitly cleared — the root Run function will not pass CompareWith

```go
func Setup(t *testing.T, req *Request) error {
	req.CompareWith = ""
	return nil
}
```
