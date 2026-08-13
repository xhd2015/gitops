# Scenario

**Feature**: `FetchTargetedRefs` with no refs is an error (does not start a full-heads fetch).

```
FetchTargetedRefs(refs=[]) -> error
```

## Preconditions

- Root Setup still injects `ReposRoot` and a remote (unused for dest inventory).

## Steps

1. Group marks empty-input validity.
2. Leaf clears `Refs`.

## Context

- Split factor: **empty vs non-empty refs**.
