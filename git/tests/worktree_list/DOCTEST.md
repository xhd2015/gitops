# Worktree List / Resolve / WorktreesOnBranch Doctests

Doc-style tests for product-neutral worktree inventory APIs in
`github.com/xhd2015/gitops/git/worktree`: `List`, `ListLinked`,
`ParseListPorcelain`, `WorktreesOnBranch`, and resolve helpers.

# DSN (Domain Specific Notion)

**Participants**

- **Main checkout** — primary worktree where `.git` is a directory; always
  appears in `List`.
- **Linked worktree** — additional checkout registered via `git worktree add`;
  `.git` is a file; returned by `ListLinked`.
- **Entry** — one porcelain row: `Path`, `Branch` (empty when detached),
  `HEAD`, `IsMain`.
- **List** — all registered worktrees (main + linked), including dead paths
  until pruned.
- **ListLinked** — linked entries only (excludes main).
- **WorktreesOnBranch** — filter of registered entries whose `Branch` equals a
  named branch; no policy / no error when `len > 1`.
- **ParseListPorcelain** — pure parser of `git worktree list --porcelain` text.
- **Resolve helpers** — `ResolveMainRepo`, `IsDead` for path identity and
  missing directories.

**Behaviors**

- Detached worktrees never match a named branch filter (`Branch == ""`).
- Dead (directory removed, not pruned) entries still appear in list / filter.
- Two worktrees may share one branch (`--force` or
  `checkout --ignore-other-worktrees`); API returns data only, never refuses.
- Paths compare after clean / symlink resolution (macOS `/var` vs `/private/var`).

## Version

0.0.2

## Decision Tree

```
worktree_list
├── list/
│   ├── single-main/              (LEAF) only main on master
│   ├── main-plus-linked/         (LEAF) List includes both; ListLinked excludes main
│   └── dead-included/            (LEAF) rm -rf linked path still listed + IsDead
├── worktrees_on_branch/
│   ├── only-main-on-master/      (LEAF) master→1; other→0
│   ├── feature-one-linked/       (LEAF) feature→1 linked path only
│   ├── two-linked-same-branch/   (LEAF) two linked on feature → len=2
│   ├── main-shares-with-linked/  (LEAF) main+linked same branch → len=2
│   ├── dead-still-on-branch/     (LEAF) dead linked still filtered by branch
│   └── detached-not-matched/     (LEAF) detached excluded from named branch
├── parse_list_porcelain/
│   ├── two-entries/              (LEAF) pure parse: path + branch strip refs/heads/
│   └── detached-empty-branch/    (LEAF) detached → Branch empty
└── resolve/
    └── main-from-linked/         (LEAF) ResolveMainRepo(linked) → main path
```

## Test Index

| # | Leaf | Description |
|---|------|-------------|
| 1 | `list/single-main` | `List` → one main entry on `master`, `IsMain=true` |
| 2 | `list/main-plus-linked` | `List` len=2; `ListLinked` len=1, not main |
| 3 | `list/dead-included` | Dead linked path still in `List`; `IsDead(path)=true` |
| 4 | `worktrees_on_branch/only-main-on-master` | `master`→1; nonexistent branch→0 |
| 5 | `worktrees_on_branch/feature-one-linked` | `feature`→1 linked only; path matches linked |
| 6 | `worktrees_on_branch/two-linked-same-branch` | Two linked on `feature` → len=2 |
| 7 | `worktrees_on_branch/main-shares-with-linked` | Main + linked both on `feature` → len=2 |
| 8 | `worktrees_on_branch/dead-still-on-branch` | Dead linked still in `WorktreesOnBranch(feature)` |
| 9 | `worktrees_on_branch/detached-not-matched` | Detached not in named-branch filter |
| 10 | `parse_list_porcelain/two-entries` | Pure porcelain → 2 entries, branch stripped |
| 11 | `parse_list_porcelain/detached-empty-branch` | `detached` marker → empty Branch |
| 12 | `resolve/main-from-linked` | `ResolveMainRepo(linked)` equals main repo |

## How to Run

```sh
cd external/gitops-master-2026-07-30
doctest vet ./git/tests/worktree_list
doctest test ./git/tests/worktree_list/...
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/gitops/git/worktree"
)

// Op selects which worktree API Run exercises.
//   "list" | "list_and_linked" | "worktrees_on_branch" |
//   "worktrees_on_branch_multi" | "parse_list" | "resolve_main" | "is_dead"
type Request struct {
	Op         string
	Dir        string // main repo path for live git APIs
	Branch     string // single branch for WorktreesOnBranch
	Branches   []string
	Porcelain  string // input for ParseListPorcelain
	Path       string // for ResolveMainRepo / IsDead
	LinkedPath  string // expected linked path (asserts)
	LinkedPath2 string // second linked path when two linked worktrees
	MainPath    string // expected main path (asserts)
	DeadPath    string // expected dead path (asserts)
}

type Response struct {
	Entries  []worktree.Entry
	Linked   []worktree.Entry
	ByBranch map[string][]worktree.Entry
	MainRepo string
	IsDead   bool
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	switch req.Op {
	case "list":
		entries, err := worktree.List(req.Dir)
		if err != nil {
			return nil, err
		}
		return &Response{Entries: entries}, nil
	case "list_and_linked":
		entries, err := worktree.List(req.Dir)
		if err != nil {
			return nil, err
		}
		linked, err := worktree.ListLinked(req.Dir)
		if err != nil {
			return nil, err
		}
		return &Response{Entries: entries, Linked: linked}, nil
	case "worktrees_on_branch":
		entries, err := worktree.WorktreesOnBranch(req.Dir, req.Branch)
		if err != nil {
			return nil, err
		}
		return &Response{Entries: entries}, nil
	case "worktrees_on_branch_multi":
		by := make(map[string][]worktree.Entry, len(req.Branches))
		for _, b := range req.Branches {
			entries, err := worktree.WorktreesOnBranch(req.Dir, b)
			if err != nil {
				return nil, err
			}
			by[b] = entries
		}
		return &Response{ByBranch: by}, nil
	case "parse_list":
		return &Response{Entries: worktree.ParseListPorcelain(req.Porcelain)}, nil
	case "resolve_main":
		main, err := worktree.ResolveMainRepo(req.Path)
		if err != nil {
			return nil, err
		}
		return &Response{MainRepo: main}, nil
	case "is_dead":
		return &Response{IsDead: worktree.IsDead(req.Path)}, nil
	default:
		t.Fatalf("unknown op %q", req.Op)
		return nil, nil
	}
}
```
