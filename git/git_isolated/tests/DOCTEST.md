# git_isolated

## Version

0.0.1

Hook-free git subprocess runner for tests and fixtures.

## DSN

- **Env()** — `GIT_CONFIG_GLOBAL=/dev/null`, `GIT_CONFIG_SYSTEM=/dev/null`, plus hook-framework disables.
- **Command / Run / Output** — git with default test identity (`test@test.com` / `Test`).
- **Init** — `git -c init.templateDir= init -b <branch>` (no template hooks).
- **MustRun / MustOutput** — testing.T helpers that fatal on error.

## Test Index

| # | Leaf | Description |
|---|------|-------------|
| 1 | `commit-without-global-config` | commit succeeds when global gitconfig has no user.email |
| 2 | `repo-local-hook-runs` | repo-local core.hooksPath still executes |

## How to Run

```sh
doctest vet ./git/git_isolated/tests
doctest test ./git/git_isolated/tests
```

```go
import (
	"testing"
)

type Request struct {
	RepoDir      string
	GlobalConfig string
	HooksDir     string
}

type Response struct{}

func Run(t *testing.T, req *Request) (*Response, error) {
	return &Response{}, nil
}
```