# Review: gitops vs go-best-practice

**Module:** `github.com/xhd2015/gitops`  
**Date:** 2026-08-06  
**Scope:** codebase structure, external-command usage, package layout, API surface (CLI-adjacent design).  
**Skill basis:** [go-best-practice](https://github.com/xhd2015) topics — primarily **`cmd-exec`**; **CLI UX / less-flags / kool-create / go:embed** assessed for applicability only.

This is a **library** that shells out to `git` (read path: `git/`, write path: `gitwrite/`, DTOs: `model/`). There is **no user-facing CLI**, no flag parser, and no embedded assets. Findings below are ordered by severity. Recommended changes are grounded in go-best-practice where applicable; library-specific Go layout notes are called out separately.

---

## Executive summary

| Area | Verdict |
| ---- | ------- |
| **cmd-exec** (external commands) | **Largest gap.** Three concurrent execution styles; bash pipelines for core logic; `Verbose` always on; partial adoption of `xgo/support/cmd`. |
| **Package layout** | Functional split (`git` / `gitwrite` / `model`) is sound; internal layout is uneven (flat mega-package + mixed subpackages, duplicate APIs). |
| **CLI / less-flags / skill-cli** | **N/A** — no binary entrypoint. Treat as non-goals unless you add a diagnostic CLI later. |
| **kool-create** | **N/A** — existing mature library; scaffolding recipes do not apply to retrofit. |
| **go:embed assets** | **N/A** — no UI/extension assets; doctest fixtures are markdown, not embed trees. |

**Top actions (if prioritizing fixes later):** (1) unify all git invocations on `xgo/support/cmd` (or one internal wrapper), (2) eliminate bash multi-command scripts for production APIs, (3) fix known correctness bugs (`ListCommits` literal `"ref"`, `Commit` `user.name`, error strings), (4) tidy package/API surface.

---

## Project shape (as inspected)

```text
gitops/
  git/           # read-oriented git operations (large flat package + subpackages)
    command.go   # bash runner (go-inspect/sh)
    fetch/       # fetch argv builder only
    worktree/    # worktree ops (also thin re-export at git/*.go)
    status/      # pure porcelain parsers
    url/         # remote URL parse/join
    git_isolated/# test isolation runner (os/exec)
    tests/       # large doctest tree
  gitwrite/      # mutating ops (add/commit/push/tag/restore)
  model/         # shared types/constants
  README.md
  go.mod         # go 1.14; deps: go-coverage, go-inspect, xgo
```

**Positive:**

- Clear **read vs write** package split (`git` vs `gitwrite`) reduces accidental mutation surface.
- **`git_isolated`** is a good idea for hermetic tests (identity via `-c`, strip global config/hooks).
- Substantial **doctest** coverage under `git/tests/` and unit tests for several paths.
- Newer code often already uses **`github.com/xhd2015/xgo/support/cmd`** (aligned with **cmd-exec**).
- `GetOnDiskChangedFiles` uses a clean functional-options pattern.

---

## Findings (by severity)

### Critical

#### C1. Three parallel git execution stacks (cmd-exec violation)

**Topic:** `cmd-exec`  
**Evidence:**

| Style | Where | Notes |
| ----- | ----- | ----- |
| Bash via `go-inspect/sh` / `go-coverage/sh` | `git/command.go`, `list_commit.go`, `find_merges.go`, `diff_commit.go`, `grep.go`, `fetch.go` | `set -e` + `cd` + multi-line shell scripts |
| Fluent `xgo/support/cmd` | Many newer files (`list_file.go`, `inspect_worktree.go`, `worktree/*`, most of `gitwrite`) | Matches go-best-practice **cmd-exec** |
| Raw `os/exec` | `compare_branch.go`, `check_ignore.go`, `get_staged_files.go`, `diff_cached.go`, `rev_parse_show_toplevel.go`, `is_inside_git.go`, `gitwrite/restore_staged.go`, `git_isolated` | No shared Dir/Env/error policy |

**Why it matters:** Callers get inconsistent error types, stderr capture, timeout, env, and logging. Windows portability and testability suffer. go-best-practice **cmd-exec** recommends a single fluent API (`cmd.Dir(...).Env(...).Output/Run`) so every invocation behaves the same.

**Recommendation:**

1. Introduce one internal helper, e.g. `git/internal/gitcmd` or unexported funcs on package `git`, wrapping `xgo/support/cmd` only:

   ```go
   // preferred shape (cmd-exec)
   out, err := cmd.Dir(dir).Output("git", args...)
   err := cmd.Dir(dir).Env(env).Run("git", args...)
   ```

2. Migrate bash `RunCommand*` and raw `exec.Command("git", ...)` call sites to that helper.
3. Keep `git_isolated` as a **test-only** variant of the same builder (inject Env + `-c user.*`), not a third production path.
4. Stop exporting `RunCommand` / `RunCommands` / `RunCommandsWithOptions` as public library API (or document them as deprecated shell escape hatches).

---

#### C2. Production logic implemented as bash scripts inside Go

**Topic:** `cmd-exec` (+ maintainability)  
**Evidence:** `doListRelativeToBase` in `git/list_commit.go` builds a large bash program (`[[ ]]`, `for`, `head -n1`, nested `git` calls). Similar patterns in `diff_commit.go`, `find_merges.go`, `GrepLines` (`git grep … \|\| true` via shell).

**Why it matters:**

- Requires **bash** (not `sh`, not Windows `cmd`); contradicts README “mostly compatible with … Windows”.
- Hard to unit-test branches of the algorithm without full shell.
- Dead code retained (`useV2 := true` still carries a large `else` “old implementation” shell block).
- Quoting depends on two different `sh` packages (`go-inspect/sh` and `go-coverage/sh` in `find_merges.go`).

**Recommendation:** Rewrite multi-step workflows as sequential Go calls through the unified cmd helper. Parse outputs in Go. Delete dead v1 shell paths once doctests pass.

---

#### C3. Library always enables shell Verbose logging

**Topic:** `cmd-exec` (Debug is opt-in)  
**Evidence:** `git/command.go`:

```go
opts.Verbose = true
opts.NeedStdOut = true
res, _, err := sh.RunBashWithOpts(commands, opts)
```

**Why it matters:** go-best-practice treats **Debug** as an explicit choice (`cmd.Debug().Run(...)`). A library that always prints every bash command pollutes consumer stdout/stderr and makes embedding in tools/CI noisy.

**Recommendation:** Default silent (like `cmd.Run` / `cmd.Output`). Gate verbose via env (`GITOPS_DEBUG=1`) or an optional `Options{Verbose bool}` — never hard-code `Verbose = true` on the shared path.

---

### High

#### H1. Correctness bugs in exported APIs

| Bug | Location | Detail |
| --- | -------- | ------ |
| Wrong argv when listing commits | `ListCommits` (`list_commit.go`) | When `beginRef == ""`, appends the **literal** `"ref"` instead of the `ref` variable → always logs the wrong revision. |
| Wrong `user.name` config | `gitwrite.Commit` | `-c user.name=` is set to **`authorMail`**, not `authorName`. Author identity for the commit message uses the correct `--author=`, but repo config identity is wrong. |
| Inverted error text | `ErrNotExists` (`rev_parse.go`) | Message is `"reference does exist"`; should be **does not exist**. |

**Recommendation:** Fix these first (small, high-value). Add regression tests for `ListCommits("", "HEAD")` and commit config identity.

---

#### H2. README Windows claim vs actual dependencies

**Topic:** `cmd-exec` / platform  
**Evidence:**

- README: “mostly compatible with both Unix and Windows”.
- `RunCommand` uses `set -e` and bash (`sh.RunBashWithOpts`).
- `git_isolated.Env` uses `GIT_CONFIG_GLOBAL=/dev/null` (Unix path; Windows needs `NUL` or a real empty file).
- Complex paths use bash `[[ ]]` and `head`.

**Recommendation:** Either (a) drop Windows claim and document “Unix + git + bash”, or (b) finish **cmd-exec** migration (no bash) and fix isolated env for Windows. Prefer (b) for a library named “ops”.

---

#### H3. Inconsistent Dir / `-C` / cwd semantics

**Evidence:**

- Most `cmd.Dir(dir)` set process working directory.
- Some use `git -C dir …` (`GetStagedFiles`, `DiffCached`, `RestoreStaged`).
- `IsInsideGit` sets `cmd.Dir` correctly; fine.
- `compare_branch` helpers set `cmd.Dir` via raw exec — OK locally, but bypass shared error formatting.

**Recommendation:** Standardize on **one** approach. Preferred: `cmd.Dir(dir).Output("git", subcmd, …)` without `-C` (matches **cmd-exec** examples). Document that `dir` is always the git worktree root or a valid cwd for git.

---

#### H4. Errors and exit codes handled inconsistently

**Evidence:**

- `diffFiles`: on `*exec.ExitError`, returns **empty string and nil error** (swallows real failures; `|| true` semantics buried in Go).
- `GrepLines`: shell `\|\| true` hides exit 1 (no match) **and** other failures if not careful.
- `IsAncesorOf` / worktree clean: ExitError → `false, nil` (reasonable for merge-base / diff --quiet).
- `checkGitInstalled` only used from a couple of entrypoints (`IsSubmodule`, `ListIgnoredDirs`); many APIs still surface opaque `exec: "git": executable file not found`.

**Recommendation:** Centralize exit-code policy per git subcommand family. Map missing git → `ErrGitNotInstalled` once at the helper layer. Never convert all ExitErrors to success.

---

### Medium

#### M1. Package layout: flat mega-package + uneven subpackages

**Layout issues:**

| Pattern | Issue |
| ------- | ----- |
| ~50 files in package `git` | Hard to navigate; unrelated concerns (blame, fetch, worktree thin wrappers, grep, merge graph). |
| `git/worktree` + `git/worktree.go` + `git/worktree_clean.go` | Thin re-exports (`AcquireTempWorkTree` → `worktree.AcquireTemp`) add dual import paths. |
| `git/fetch` only formats argv | Could be private helpers in `fetch.go`; subpackage overkill unless shared outside module. |
| `git/status` pure parsers | Good separation, but porcelain classification is **duplicated** in `inspect_worktree.go` with different rules. |
| `package git_isolated` | Underscore package name is non-idiomatic Go (`gitisolated` / `isolated` preferred). |
| `find_merges.go` aliases `model` as `git` | `import git "github.com/xhd2015/gitops/model"` confuses readers (`[]*git.MergePoint` is model, not package git). |

**Recommendation (library layout, not a go-best-practice topic):**

- Keep public entry packages shallow: `git`, `gitwrite`, `model`.
- Move implementation into `git/internal/...` if you need folders without expanding the import graph.
- Prefer **one** porcelain parser package; have `InspectWorktree` call it.
- Rename `git_isolated` → `gitisolated` (or nest under `git/internal/isolated`) on next major.

---

#### M2. Overlapping / confusing public APIs

| Cluster | Functions | Problem |
| ------- | --------- | ------- |
| Branches containing ref | `GetBranchesContainingRef`, `SearchBranchesContainingRef`, `GetBranchesHoldingRef`, `SearchRefsContainingRef` | Similar names, different filters (first-parent vs not); hard to choose. |
| Clean checks | `WorkTreeClean`, `IndexClean` (package git) vs `IsClean`, `IndexClean`, `IsDiffClean`, `IsPorcelainClean` (worktree) | Dual surface; `IsClean` vs `IsPorcelainClean` semantics differ (untracked). |
| Options style | `*FetchOptions`, `*ListBranchRefOptions`, `GetCommitsOptions` variadic, functional options for on-disk files | Inconsistent extension story for consumers. |

**Recommendation:** Document a single preferred function per cluster; deprecate the rest. Prefer one options style for new APIs (functional options **or** opts structs — pick one). Align with “less flags / clear surface” spirit of **flags-parsing** even though this is a library: fewer knobs, better defaults.

---

#### M3. Dual `sh` dependencies (`go-inspect` + `go-coverage`)

**Evidence:** Most bash quoting uses `go-inspect/sh`; `find_merges.go` imports `go-coverage/sh`. `go.mod` requires both. go-coverage is a heavy transitive surface for a git helper library.

**Recommendation:** After bash removal, drop both if only used for Quote/RunBash. Prefer Go’s own quoting avoidance (argv slices, never shell). Until then, use **one** Quote helper.

---

#### M4. Stale Go version in `go.mod` (`go 1.14`)

**Implications:**

- Modern module consumers default much higher; `go 1.14` signals unmaintained.
- **go:embed** (skill topic) needs ≥ 1.16 — not used here, but version blocks modern stdlib.
- Tooling and dependency resolution behave better with a current language version (e.g. 1.22+).

**Recommendation:** Bump `go` directive after CI verification; not a feature change.

---

#### M5. No `context.Context` on long operations

Fetch supports timeout only via bash `RunBashOptions.Timeout`. No API accepts `context.Context` for cancel/deadline.

**Recommendation:** When consolidating on **cmd-exec**, add `Context` to the internal runner if `xgo/support/cmd` supports it, or document that cancellation is unsupported. At minimum, plumb timeout consistently without bash.

---

#### M6. Test / production API bleed

**Evidence:** `TestExported_ClassifyPorcelainLine` is an exported production symbol “for doctests”.

**Recommendation:** Prefer doctest hooks in `_test.go`, `export_test.go`, or a `git/internal/porcelain` package. Do not grow `TestExported_*` on the public surface.

---

#### M7. `gitwrite` quality uneven

| Good | Weak |
| ---- | ---- |
| Mostly uses `cmd.Dir` (**cmd-exec**) | `RestoreStaged` uses raw `os/exec` |
| Proxy env helper for push | `Commit` user.name bug (H1) |
| Separates mutation from `git` | Little validation on `AddAll` / empty paths |

**Recommendation:** Align `RestoreStaged` with `cmd`; fix Commit; consider optional dry-run only if a CLI is added later (**cli/dry-run** is N/A for pure library today).

---

### Low

#### L1. Naming and typos (exported)

| Symbol / text | Issue |
| ------------- | ----- |
| `IsAncesorOf` | Missing ‘t’ — **exported** typo; needs deprecation alias `IsAncestorOf` |
| `Unshawllow` | Typo for Unshallow |
| `ErrNotExists` text | See H1 |
| Comments / README | `referes`, `pseduo`, `requries`, `intial purepose` |
| Constants | `COMMIT_WORKING`, `ORIGIN`, `MASTER` use SCREAMING_SNAKE; Go prefers `CommitWorking` / mixed caps for exported consts |

---

#### L2. Sparse package docs

- No `// Package git ...` doc comments.
- README is four sentences; no install, no examples, no “read vs write”, no Windows caveat.

**Recommendation:** Expand README with module install, minimal example using `cmd`-style APIs, and package boundaries. Package comments help godoc consumers.

---

#### L3. Topics assessed as not applicable

| Topic | Why N/A |
| ----- | ------- |
| **flags-parsing / less-flags** | No CLI argv. If you later add `gitops` debug CLI, use less-flags + per-subcommand `--help` (`flags-parsing/subcommand`). |
| **cli/color, streaming, dry-run, skill-cli, inline-tui-mouse** | No TUI/CLI host. Library should stay silent by default (ties back to C3). |
| **kool-create** | For greenfield scaffolds (`server`, `go-react`, …), not retrofitting this module. |
| **go-embed-assets** | No generated frontend/extension trees. Doctest markdown is not an embed/hydrate case. |

Optional future: a thin **`cmd/gitops`** diagnostic tool (status/inspect) would be the first place to apply less-flags + cmd-exec + streaming — keep it out of the library packages.

---

## Recommended change plan (when you implement — not done in this review)

Phased so each step is independently verifiable (doctests + `go test ./...`).

### Phase 1 — Safety / correctness (small)

1. Fix `ListCommits` literal `"ref"` bug; add unit/doctest.
2. Fix `gitwrite.Commit` `user.name=` to use `authorName`.
3. Fix `ErrNotExists` message.
4. Default `Verbose` off in `RunCommandWithOptions`.

### Phase 2 — cmd-exec unification

1. Add internal `runGit(dir, args...)` / `outputGit(dir, args...)` wrapping `xgo/support/cmd`.
2. Migrate raw `os/exec` git call sites and `gitwrite.RestoreStaged`.
3. Wire `ErrGitNotInstalled` once in the helper.
4. Point `git_isolated` at shared builder + test Env.

### Phase 3 — Remove bash production paths

1. Port `GetCommit` / `GetCommits` / `ListCommitRelativeToBase` / `DiffCommit` / `FindMerge*` / `GrepLines` / `Fetch*` to argv-based Go.
2. Delete `RunCommand*` public API or move behind `// Deprecated`.
3. Drop `go-inspect` / `go-coverage` if unused; tidy `go.mod`.
4. Bump `go` directive; re-verify CI/doctests on Linux + macOS; decide Windows stance.

### Phase 4 — API / layout cleanup

1. Deprecate duplicate branch/clean helpers; document preferred names.
2. Single porcelain classifier; remove `TestExported_*`.
3. Collapse thin re-exports or document “import `git/worktree` for full API”.
4. Rename `IsAncesorOf` → `IsAncestorOf` with compatibility wrapper.
5. Expand README + package docs.

---

## Mapping to go-best-practice topics

| Finding | Topic |
| ------- | ----- |
| C1, C2, C3, H2, H3, H4, M3, M5, M7 | **`cmd-exec`** |
| M2 (fewer knobs, clear surface) | Spirit of **`flags-parsing`** / “less flags” (library analogue) |
| Silent default library I/O | Spirit of **`cli`** (no accidental streaming/noise); **`cli/streaming`** only if a CLI is added |
| Side-effect writes | **`cli/dry-run`** only if a future CLI wraps `gitwrite` |
| Scaffolding | **`kool-create`** — N/A |
| Assets | **`go-embed-assets`** — N/A |

---

## What is already in good shape

1. **Intentional package split** between read (`git`) and write (`gitwrite`).
2. **Partial cmd-exec adoption** — many hot paths already use `xgo/support/cmd` with `.Dir()`.
3. **`git_isolated`** design for tests (identity + ignore global config/hooks).
4. **Rich behavioral tests** (doctests + table-style unit tests) for merge graphs, worktree, on-disk changes.
5. **Typed models** with JSON tags for tooling integration.
6. **Functional options** on `GetOnDiskChangedFiles` as a template for new APIs.

---

## Out of scope for this review

- Implementing the fixes above (per request: review only; no non-trivial code changes).
- Performance profiling of git subprocess fan-out.
- Full security review of URL/token helpers in `git/url` (token-in-URL patterns deserve a separate pass if used with secrets).

---

*Generated against worktree branch `master-2026-08-06-use-go-best-practice-to-review-current-project` using the go-best-practice skill index and topic recipes (`cmd-exec`, `cli`, `flags-parsing`, `kool-create`, `go-embed-assets`).*
