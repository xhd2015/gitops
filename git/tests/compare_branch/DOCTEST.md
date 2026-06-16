# CompareBranches Test Case Tree

Run with:
```sh
doctest test ./ -v
```

## Text Tree

```
CompareBranches(dir, refA, refB)
├── 🍃 identical-leaf
├── 🍃 fast-forward-a-to-b
├── 🍃 fast-forward-b-to-a
├── 🍃 divergent-leaf
├── 🍃 divergent-no-file-diff-leaf
├── 🔴 errors-invalid-ref-a
├── 🔴 errors-invalid-ref-b
├── 🔴 not-a-git-repo-leaf
├── 🍃 dir-param-valid
└── 🔴 dir-param-nonexistent
```

## Test Case Index

| # | Path | Preconditions | Expected |
|---|---|---|---|
| 1 | `identical-leaf` | Both refs resolve to same commit | `BranchRelationSame` |
| 2 | `fast-forward-a-to-b` | `a` is ancestor of `b` — `b` can ff | `BranchRelationAIsAncestorOfB`, `CommitsAheadB=1` |
| 3 | `fast-forward-b-to-a` | `b` is ancestor of `a` — `a` can ff | `BranchRelationBIsAncestorOfA`, `CommitsAheadA=1` |
| 4 | `divergent-leaf` | Both have unique commits, file diff | `BranchRelationDiverged`, `DiffFileCount>0`, merge base non-empty |
| 5 | `divergent-no-file-diff-leaf` | Both have unique commits, same file content | `BranchRelationDiverged`, `DiffFileCount=0` |
| 6 | `errors-invalid-ref-a` | `refA` does not resolve | error with ref name |
| 7 | `errors-invalid-ref-b` | `refB` does not resolve | error with ref name |
| 8 | `not-a-git-repo-leaf` | `dir` is not a git repo | error |
| 9 | `dir-param-valid` | `dir` points to a valid git repo | `BranchRelationSame` |
| 10 | `dir-param-nonexistent` | `dir` does not exist | error |
