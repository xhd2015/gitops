# Targeted dest refs (prefix + hash + conservative fetch)

`github.com/xhd2015/gitops/repos`:

- `FetchTargetedRefs` writes dest refs as `{prefix}/{sha256hex(input)}`.
- Default prefix is `refs/gitops/targets`.
- `TargetRefPrefix` customizes dest prefix; trailing slashes do not produce `//`.
- Fetch is conservative: `--filter=blob:none`, no `+refs/heads/*` mirror.

**Layer: L2** — in-process library call against a local `file://` remote.
Leaves may `t.Parallel()`. No network, no e2e label.

## Version

0.0.2

# DSN (Domain Specific Notion)

Participants:

- **Caller** invokes `repos.FetchTargetedRefs(ctx, cloneURL, refs, opts)` with
  injectable `ReposRoot`, optional `TargetRefPrefix`, and a `Progress` writer.
- **Targeted dest ref** is `{prefix}/{sha256hex}` where the suffix is lowercase
  hex of SHA-256 over the **exact caller input** (not the remapped source).
- **Default prefix** is `refs/gitops/targets`. Empty `TargetRefPrefix` means default.
- **Custom prefix** (e.g. `refs/other/targets`) uses the same hash rule.
- **Slash normalize** — prefix with or without a trailing slash yields one
  `{prefix}/{hash}` form (no double slash).
- **Conservative fetch** — argv matches
  `git fetch --no-tags --filter=blob:none --depth=1 origin +<source>:<dest>`
  for requested refs only; not a full-heads mirror (`+refs/heads/*`).
- **Partial cache** lives under `{ReposRoot}/targeted-cache/…` (no second cache root).
- **Local file:// remote** has two branches (`master`, `feature`). No GitLab.
- **Harness** places `ReposRoot` and the remote under `d.DOCTEST_CASE`. No
  `os.Setenv` / `t.Setenv` / `os.Chdir` / `t.Chdir`.

Behaviors:

1. Fetch one branch with default prefix → dest is `refs/gitops/targets/<64 hex>`;
   no `refs/other/targets/…` dest.
2. `TargetRefPrefix=refs/other/targets` → dest uses that prefix + the same hash.
3. Prefix `refs/gitops/targets/` (trailing slash) → still `refs/gitops/targets/<hash>`.
4. Observed fetch argv contains `--filter=blob:none` and does not contain `+refs/heads/*`.
5. Remote has `master` + `feature`; fetch `feature` only → feature dest exists;
   unused branch is not created as a dest target.
6. `FetchTargetedRefs` with no refs returns an error.

## Decision Tree

```text
targeted-refs
|-- dest-prefix                      # dest naming: prefix + hash identity
|   |-- default                      # empty TargetRefPrefix → refs/gitops/targets/<hash>
|   |-- custom                       # TargetRefPrefix honored; same hash rule
|   `-- slash-normalize              # trailing slash → no double slash
|-- conservative-fetch               # not a full mirror
|   |-- no-full-heads                # argv has --filter=blob:none; no +refs/heads/*
|   `-- only-requested               # feature dest exists; master dest absent
`-- empty-refs                       # input validity
    `-- error                        # no refs → error
```

### Parameter significance (high → low)

1. **Input validity** — empty refs (error) vs non-empty (fetch)
2. **Dest prefix policy** — default / custom / trailing-slash normalize
3. **Conservative fetch shape** — argv vs dest inventory for unrequested refs

## Test Index

| Leaf | Fixture | Expect |
|------|---------|--------|
| `dest-prefix/default` | fetch `master`, empty prefix | dest `refs/gitops/targets/<sha256(master)>`; no `refs/other/targets/…` |
| `dest-prefix/custom` | fetch `master`, prefix `refs/other/targets` | dest that prefix + same hash; no default-prefix dest |
| `dest-prefix/slash-normalize` | fetch `master`, prefix `refs/gitops/targets/` | dest `refs/gitops/targets/<hash>` (not `…/targets//<hash>`) |
| `conservative-fetch/no-full-heads` | fetch `master`, Progress captured | Progress contains `--filter=blob:none`; does not contain `+refs/heads/*` |
| `conservative-fetch/only-requested` | remote `master`+`feature`; fetch `feature` | feature dest exists; master dest not created |
| `empty-refs/error` | `refs` empty / nil | `FetchTargetedRefs` returns an error |

## How to Run

From the gitops module root:

```sh
doctest vet ./repos/tests/targeted-refs
doctest test -v ./repos/tests/targeted-refs
doctest test -v ./repos/tests/targeted-refs/dest-prefix/default
```

No network / e2e labels — local `file://` L2.

### Target API surface (implementer owns; designer does not write product)

Package path:

```text
github.com/xhd2015/gitops/repos
```

```go
package repos

import (
	"context"
	"io"
	"time"
)

const DefaultTargetRefPrefix = "refs/gitops/targets"

// TargetedFetchOptions controls a narrow fetch into a partial bare cache.
type TargetedFetchOptions struct {
	Auth            *GitAuthConfig
	Env             map[string]string
	MaxAge          time.Duration
	Depth           int
	Deepen          int
	Unshallow       bool
	HistoryPhase    string
	ReposRoot       string    // empty → ReposRoot() (~/.repos); tests always inject
	Progress        io.Writer // prints fetch argv (dir: + $ git …)
	TargetRefPrefix string    // empty → DefaultTargetRefPrefix; trim trailing /
}

type TargetedFetchResult struct {
	CacheDir       string
	Commits        map[string]string
	Fetched        bool
	FilterAccepted bool
	LockWait       time.Duration
	FetchDuration  time.Duration
}

// FetchTargetedRefs fetches only the requested refs into {ReposRoot}/targeted-cache/….
// Dest local ref is {normalizedPrefix}/{hex(sha256(input))}.
func FetchTargetedRefs(ctx context.Context, cloneURL string, refs []string, opts TargetedFetchOptions) (*TargetedFetchResult, error)
```

Optional L1 helper (not required by these leaves): `TargetDestRef(prefix, input string) string`.

### Out of scope

- Changing `git.FetchAll` / `FetchSingle`
- Bumping gitops Go version (keep 1.14)

```go
import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/gitops/git/git_isolated"
	"github.com/xhd2015/gitops/repos"
)

const (
	DefaultTargetRefPrefix = "refs/gitops/targets"
	CustomTargetRefPrefix  = "refs/other/targets"
)

type Request struct {
	CloneURL        string
	Refs            []string
	ReposRoot       string
	TargetRefPrefix string
	Depth           int
}

type Response struct {
	CacheDir string
	Commits  map[string]string
	Fetched  bool
	Progress string
	AllRefs  []string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	depth := req.Depth
	if depth <= 0 {
		depth = 1
	}
	var progress bytes.Buffer
	result, err := repos.FetchTargetedRefs(context.Background(), req.CloneURL, req.Refs, repos.TargetedFetchOptions{
		ReposRoot:       req.ReposRoot,
		TargetRefPrefix: req.TargetRefPrefix,
		Depth:           depth,
		Progress:        &progress,
	})
	resp := &Response{Progress: progress.String()}
	if result != nil {
		resp.CacheDir = result.CacheDir
		resp.Commits = result.Commits
		resp.Fetched = result.Fetched
		resp.AllRefs = listRefNames(result.CacheDir)
	}
	return resp, err
}

func listRefNames(cacheDir string) []string {
	if cacheDir == "" {
		return nil
	}
	out, err := git_isolated.Output(cacheDir, "for-each-ref", "--format=%(refname)")
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	lines := strings.Split(out, "\n")
	refs := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			refs = append(refs, line)
		}
	}
	return refs
}
```
