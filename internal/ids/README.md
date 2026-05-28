# internal/ids

The central allocation of stable integer identifiers for every card, hero, weapon, weapon
ability, and token in the build. Keeping the IDs in one package lets per-entity caches (attack-turn
step text, display name, attacker meta) be plain slices indexed by ID instead of map lookups
on the hot path.

## Key types

- `CardID` (`uint16`) — identifies a printed card. Each pitch variant (Red / Yellow / Blue) is
  a distinct card with its own ID. `InvalidCard` (0) is the reserved zero sentinel; valid IDs
  start at 1.
- `HeroID` (`uint16`) — identifies a hero printing. `InvalidHero` (0) is the zero sentinel.
  Same width as `CardID` so (hero, card) tuples stay fixed-size integer structs.

Weapons are cards: the platonic weapon card and its activated ability both carry `CardID`s
in the shared number space (see "ID anchoring scheme"). There is no separate weapon-ID type.

## ID anchoring scheme

All card-shaped IDs share one `CardID` number space, carved into anchored ranges so that
adding entries to one range never renumbers another:

1. Card IDs (`card_ids.go`) start at 1, ordered alphabetically by card name with Red then
   Yellow then Blue within each family.
2. Weapon-permanent IDs (`weapon_ids.go`) anchor immediately past the last card ID
   (`ZealousBeltingBlue + iota + 1`).
3. Weapon-ability IDs (`weapon_ids.go`) anchor past the last weapon-permanent ID — the attack turn
   runner enqueues the ability, so the cache keys off the ability ID.
4. Token-ability IDs (`token_ids.go`) anchor past the last weapon-ability ID.
5. Token aura / item IDs (`token_ids.go`) anchor past the last token-ability ID.
6. Test-fake card IDs (in `internal/testutils`) anchor past the last token ID.

Each range begins with `<First> = <PreviousRangeLast> + iota + 1` so the ranges stay contiguous
and non-overlapping.

## Important files

- `card_ids.go` — `CardID`, `InvalidCard`, every card variant ID.
- `hero_ids.go` — `HeroID`, `InvalidHero`, hero IDs (`DefaultHeroID`, `ViseraiID`, …).
- `weapon_ids.go` — weapon-permanent card IDs and weapon-ability IDs.
- `token_ids.go` — token-ability IDs and token aura / item IDs.

## Gotchas and invariants

- IDs are stable within a build but are NOT a persistence format. Adding or removing entries
  renumbers later IDs in the same range. Treat them as opaque in-process handles; never
  serialize them.
- Adding a card means adding its variant IDs to `card_ids.go` (alphabetical, R/Y/B) before
  running the card generator — the generated registry references `ids.<VariantID>`.
- The zero value of every ID type is the reserved `Invalid*` sentinel; never assign it to a
  real entity.
