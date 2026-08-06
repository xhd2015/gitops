# CI workflow note

## Status

**added** (and pushed)

## Branch / remote / push

| Field | Value |
| ----- | ----- |
| Branch | `master-2026-08-06-use-go-best-practice-to-review-current-project` |
| Remote | `origin` → `ssh://git@github.com/xhd2015/gitops` |
| Push result | **success** — `1e13909..5d5dc0c` → `origin/master-2026-08-06-use-go-best-practice-to-review-current-project` (upstream set) |
| Commit SHA | `5d5dc0c55029ab51fc65235b92022e2f28501693` (`5d5dc0c`) |
| Commit message | `ci: add GitHub Actions test workflow aligned with doctest pattern` |

## Paths changed

| Path | Action |
| ---- | ------ |
| `.github/workflows/test.yml` | **added** — full Test job (go test, doctest discovery + e2e, xgo merge, summary, artifacts) |
| `script/ci/coverage-package-table.py` | **added** — package coverage markdown table for `github.com/xhd2015/gitops/` |
| `CI-WORKFLOW-NOTE.md` | **added** — this note (may land in a follow-up commit on the same branch) |

Not included in the CI commit: local review artifacts (`REVIEW-*.md`).

## How to view Actions for this push

- Repo: https://github.com/xhd2015/gitops  
- Actions (all): https://github.com/xhd2015/gitops/actions  
- Workflow runs for this branch: https://github.com/xhd2015/gitops/actions?query=branch%3Amaster-2026-08-06-use-go-best-practice-to-review-current-project  
- Workflow file on branch: https://github.com/xhd2015/gitops/blob/master-2026-08-06-use-go-best-practice-to-review-current-project/.github/workflows/test.yml  

Look for workflow name **Test** on the latest push to that branch.

## How this differs from doctest’s workflow

Reference: doctest worktree `.github/workflows/test.yml`.

| Aspect | doctest | gitops (this repo) |
| ------ | ------- | ------------------ |
| Module / `COVERPKG` | `github.com/xhd2015/doctest/...` | `github.com/xhd2015/gitops/...` |
| setup-go | `go-version-file: go.mod` | `go-version: "1.25.x"` — `go.mod` still says `go 1.14`; modern toolchain required for `doctest@latest` |
| Install doctest | `go install ./cmd/doctest` (this checkout) | `go install github.com/xhd2015/doctest/cmd/doctest@latest` |
| Doctest discovery | `--label '!e2e'` | Same |
| Doctest e2e | `--label e2e` | Same stage kept for pattern parity (repo currently has no e2e-labeled leaves) |
| Coverage merge | gotest + discovery + e2e | Same three profiles when present |
| Package table | `script/ci/coverage-package-table.py` (doctest module filter) | Same helper path; filters `github.com/xhd2015/gitops/` |
| Artifacts | coverage profiles + func report | Same set of paths |

## Note on local pre-commit hook

`git-hook-github-workflow-test` warns that this file differs from its **simpler** recommended template (`doctest test -v --label-all`). That is intentional: this workflow follows the fuller **doctest project** CI pattern (coverage, staged labels, merge, summary, artifacts) as requested.
