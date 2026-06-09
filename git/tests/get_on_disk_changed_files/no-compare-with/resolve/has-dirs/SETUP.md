## Preconditions
- There is at least one untracked directory on disk
- git status --porcelain reports directory entries (e.g. `?? view/`)

## Steps
1. Create the untracked directory structure (specified by leaf SETUP.md)
2. Verify git status returns directory paths before expansion
