# Status Porcelain Parse Doctests

Doc-style tests for product-neutral porcelain parsing in
`github.com/xhd2015/gitops/git/status`. Parse only — no wrk format strings,
no `*Wrk*` type names.

# DSN (Domain Specific Notion)

**Participants**

- **Caller** — supplies `git status --porcelain` text (no subprocess in pure
  leaves).
- **Counts** — backup-style multi-bucket aggregate: Modified, Added, Deleted,
  Untracked, Renamed, Copied, Unmerged.
- **ChangeCounts** — four-bucket aggregate: Added, Changed, Renamed, Deleted
  (inspect-compatible taxonomy; untracked `??` counts as **Added**).
- **ParsePorcelain** — maps porcelain lines → `Counts`.
- **ParseChangeCounts** — maps porcelain lines → `ChangeCounts`.

**Behaviors**

- Empty porcelain → all zeros for both parsers.
- Backup: `??` → Untracked; `M` → Modified; `A` → Added; `D` → Deleted;
  `R` → Renamed; `C` → Copied; `U` → Unmerged.
- Four-bucket: `??` → Added; `R` → Renamed; `D` → Deleted; `A` → Added;
  `M` → Changed (either XY column; first match order as inspect classifier).
- No `FormatWrk` / `dirty (...)` product strings in this package surface.

## Version

0.0.2

## Decision Tree

```
status_parse
├── counts/
│   ├── empty/              (LEAF) empty porcelain → zero Counts
│   └── mixed/              (LEAF) M + ?? → Modified/Untracked
└── change_counts/
    ├── empty/              (LEAF) empty porcelain → zero ChangeCounts
    └── representative/     (LEAF) ??/M/R/D → Added/Changed/Renamed/Deleted
```

## Test Index

| # | Leaf | Description |
|---|------|-------------|
| 1 | `counts/empty` | Empty string → all `Counts` fields zero |
| 2 | `counts/mixed` | Two ` M` + one `??` → Modified=2, Untracked=1 |
| 3 | `change_counts/empty` | Empty string → all `ChangeCounts` fields zero |
| 4 | `change_counts/representative` | `??`→Added, `M`→Changed, `R`→Renamed, `D`→Deleted |

## How to Run

```sh
cd external/gitops-master-2026-07-30
doctest vet ./git/tests/status_parse
doctest test ./git/tests/status_parse/...
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/gitops/git/status"
)

// Op: "counts" | "change_counts"
type Request struct {
	Op        string
	Porcelain string
}

type Response struct {
	Counts       status.Counts
	ChangeCounts status.ChangeCounts
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	switch req.Op {
	case "counts":
		return &Response{Counts: status.ParsePorcelain(req.Porcelain)}, nil
	case "change_counts":
		return &Response{ChangeCounts: status.ParseChangeCounts(req.Porcelain)}, nil
	default:
		t.Fatalf("unknown op %q", req.Op)
		return nil, nil
	}
}
```
