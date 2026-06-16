# IsInsideGit Test Case Tree

Run with:
```sh
doctest test ./ -v
```

## Test Case Index

| # | Path | Preconditions | Expected |
|---|---|---|---|
| 1 | `inside-git/` | `dir` is a git repo (initialised) | `(true, nil)` |
| 2 | `outside-git/` | `dir` is NOT a git repo | `(false, nil)` |
