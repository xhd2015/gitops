## Steps
1. Write `mod.go` with content "original", add and commit
2. Write `del.go` with content "original", add and commit
3. Write `ren.go` with content "original", add and commit
4. Modify `mod.go` to "modified"
5. Write `added.go` with content "new staged" and `git add added.go`
6. Write `untracked.go` with content "untracked"
7. Delete `del.go` from disk
8. Run `git mv ren.go ren_new.go`
