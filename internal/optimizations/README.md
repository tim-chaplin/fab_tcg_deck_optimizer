# internal/optimizations

## Purpose

Hot-path performance caches that sit beside the simulator without changing its behaviour.
Today the package holds one cache: a memoised chain-step text table that eliminates the
per-`Play` string concatenation and `DisplayName` allocation the chain runner would
otherwise pay on every card.

## Key types

The package exposes no types — only the `WarmChainStepCache` function. The cache itself
(`chainStepCache`, an array of `atomic.Pointer[string]`) is package-private.

## How it works

`gameengine.ChainStepText` is a swappable function variable. The package's `init` captures
the original (`bareChainStepText`) and installs `cachedChainStepText` over it. The
chain-step text depends only on `(Card.ID, FromArsenal)` — display name, types, and verb
selection are all static — so results live in a pre-warmed table indexed by a `uint32`
that packs the card ID into bits 0-15 and the `FromArsenal` flag into bit 16 (two adjacent
rows per card).

## How to use / extend

- `WarmChainStepCache(cards)` populates both the in-hand and from-arsenal entries for every
  non-nil card. It is idempotent. The registry package's `init`, and `cmd/fabsim`'s `init`,
  call it once with the full registry slice so the runtime hot path is pure cache reads.
- A card created outside the registry (a test fake or stub) still works: the cache miss
  path (`chainStepTextSlow`) computes and backfills the entry on first sighting.
- To add another hot-path cache, follow the same shape: a swappable function variable in
  the producing package, an `init` here that installs the memoised version, and a `Warm…`
  entry point callers invoke once at startup.

## Important files

- `chain_step_cache.go` — the entire cache: the `init` swap, the index packing, the warm
  and slow paths.

## Gotchas

- The cache is sized for the full `uint16` ID space (`1 << 17` entries) so lookups are
  direct bounds-checked array reads with no map or hashing.
- Multiple goroutines computing the same entry race-safely converge: every writer produces
  the same string, so a read after a race still matches spec.
- If chain-step text ever starts depending on something beyond `(Card.ID, FromArsenal)`,
  this cache becomes incorrect — the keying assumption must be revisited.
