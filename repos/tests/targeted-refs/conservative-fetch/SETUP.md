# Scenario

**Feature**: targeted fetch is conservative — blob-less, requested refs only, not a full-heads mirror.

```
FetchTargetedRefs(requested)
  -> git fetch --filter=blob:none origin +<source>:<dest>
  -> no +refs/heads/*; unused dest not created
```

## Preconditions

- Root Setup provides a two-branch (`master`, `feature`) `file://` remote.
- Progress writer is injected by root `Run`.

## Steps

1. Group marks conservative-fetch concern (argv vs dest inventory).
2. Leaves fetch a subset of remote heads and inspect Progress / dest refs.

## Context

- Split factor: **fetch shape** (argv vs which dest refs appear).
