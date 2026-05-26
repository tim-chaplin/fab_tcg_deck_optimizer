# internal/deck

## Purpose

A candidate Flesh and Blood deck — hero, equipped weapons, the 40-card deck list, plus the
user's sideboard and equipment lists — together with the mutation enumeration the anneal
search drives over and the per-deck evaluation result types. The package depends only on
narrow `Hero` / `Weapon` / `Card` / `Registry` contracts declared locally, so any concrete
card / weapon / hero / registry implementation can be swapped in; the simulator and game
engine live elsewhere.

## Key types

- `Deck` — hero, `Weapons`, an unexported `cards` slice that doubles as the runtime deck,
  plus user-managed `Sideboard` and `Equipment` name lists the simulator never reads.
- `Mutation` — one candidate single-slot or pair change: a mutated `*Deck` plus a
  human-readable `Description`.
- `CardPair` / `CardGroup` — synergy pairings (e.g. Sun Kiss / Moon Wish) the mutation
  generator adds as an atomic 2-for-2 swap. Registered in the `CardPairs` slice.
- `Stats`, `CycleStats`, `BestTurn`, `CardMarginalStats` — the aggregate hand-value
  statistics a simulation run produces. Pure data with small derived-stat methods and no
  sim / gameengine dependency, so this file sits below those packages without a cycle.
- `Defaults` / `SideboardDefault` — the equipment + sideboard loadout `ApplyDefaults` tops
  a saved deck up toward. `ViseraiDefaults` is the concrete Viserai loadout.

## How to use / extend

- Build a deck with `New(hero, weapons, cards)` — it panics if the weapon loadout breaks
  the "0-2 weapons; if 2, both 1H" rule.
- `Random(hero, size, maxCopies, rng, registry)` generates a fresh legal deck for search
  starting points.
- `AllMutations(d, maxCopies, includePairs, registry)` returns every weapon-loadout,
  single-card-swap, and (when `includePairs`) synergy-pair mutation in a deterministic
  ID-sorted order (the anneal driver shuffles it each round to keep exploration unbiased).
- The runtime methods `Shuffle` / `Draw` / `PeekTop` / `PutTop` / `PutBottom` / `Tutor`
  mutate `cards` directly; callers running an evaluation trial must `Copy()` the master
  deck first. `ShallowCopy` / `ShallowCopyFrom` / `CopyFrom` are allocation-light variants
  for the hot path.
- To register a new synergy pair, add `CardGroup` vars and append a `CardPair` to
  `CardPairs`.

## Important files

- `deck.go` — the `Deck` type, copy variants, runtime deck methods, `New`, `Random`,
  `ApplyDefaults`.
- `mutate.go` — `AllMutations` and the single-slot / pair / weapon-loadout generators.
- `stats.go` — `Stats` and the derived per-deck / per-card result types.
- `interfaces.go` — the narrow `Hero` / `Weapon` / `Card` / `Registry` contracts.
- `format.go` — composition-level formatting helpers (`NameCounts`, `DisplayNames`,
  `PitchCounts`).
- `weapon_loadouts.go` — weapon validation and legal-loadout enumeration.
- `viserai_defaults.go` — `ViseraiDefaults`.

## Gotchas

- `cards` is unexported deliberately: external code can't peek the runtime order. Inspect
  composition through `UniqueIDs` / `NameCounts` / `DisplayNames` / `PitchCounts`.
- `Shuffle` panics on a `ShallowCopy`-produced wrapper (`mustNotShuffle` is set): those
  share slice backing with peer wrappers and an in-place shuffle would corrupt them. A card
  that shuffles the deck mid-attack-turn trips this — if that becomes intentional, the
  per-permutation call site must revert to a deep `Copy`.
- `Draw` / `PeekTopN` return slices that alias the deck's backing storage; do not retain or
  mutate them past the next deck mutation.
- The mutation generators emit cap-blind candidates; `maxCopies` is enforced once by
  `filterMaxCopiesViolations` as a final post-pass in `AllMutations`.
