## Steps
1. Make at least one additional commit so HEAD has a parent
2. Get HEAD~1 commit hash: `git rev-parse HEAD~1`
3. Set req.CompareWith to that commit hash
4. Write `a.go` with content "package main", add and commit
5. Working tree is clean
