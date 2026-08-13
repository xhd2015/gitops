# Scenario

**Feature**: dest ref name is `{prefix}/{sha256hex(input)}` with a default, custom, or slash-normalized prefix.

```
FetchTargetedRefs(refs, TargetRefPrefix)
  -> dest {normalizedPrefix}/{sha256(input)}
  -> no other prefix family created
```

## Preconditions

- Root Setup provides injectable `ReposRoot` and a two-branch `file://` remote.
- Default request fetches `master`.

## Steps

1. Group marks dest-prefix concern (default / custom / slash-normalize).
2. Leaves set `TargetRefPrefix` and assert dest identity.

## Context

- Split factor: **prefix policy**.
- Hash identity is the SHA-256 of the exact caller input (`master`).
