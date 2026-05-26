# internal/registry

The master roster of every implemented card, weapon, and hero — the deck builder's source of
truth for what exists and what is legal to put in a deck. It provides ID lookups, name
lookups, and the filtered deck-construction pools.

## Key types

- `Registry` — the deck-construction view. Zero-sized; `LegalCards` / `LegalWeapons` return
  the pools with `NotImplemented` / `Unplayable` / `NotSilverAgeLegal` entries filtered out.
  Satisfies `internal/deck.Registry`; callers pass `Registry{}` to `deck.Random`.
- `Card` / `Hero` / `Weapon` (`registry.go`) — minimal local interfaces (identity + display
  name, and `Hands` for weapons). The package declares its own narrow interfaces so its
  surface stays decoupled from the sim's richer contracts.
- `NotImplemented` / `Unplayable` / `NotSilverAgeLegal` (`registry.go`) — marker interfaces.
  `isExcludedFromPool` matches them structurally to gate cards and weapons out of the
  construction pools. A `NotImplemented` card is still valid in a pre-built hand (it evaluates
  on its static stats); the markers only stop the optimizer from introducing the card.

## Registry / sim split

`internal/registry` declares minimal `Card` / `Hero` / `Weapon` interfaces (identity + display
name) so its surface stays decoupled from the sim's richer contracts. Concrete card / weapon /
hero types satisfy both, and callers needing behaviour assert to `card.Card` / `sim.Weapon` /
`sim.Hero` at the read site. The marker interfaces (`NotImplemented`, `Unplayable`,
`NotSilverAgeLegal`) are declared locally in both the registry and the sim and matched
structurally — neither package imports the other.

Callers that need cards / weapons / heroes import `internal/registry` directly and assert to
`card.Card` / `sim.Weapon` / `sim.Hero` at the read site. `gameengine.AttackStepText`
memoises results on `(Card.ID, FromArsenal)` and lazily backfills on the first call per
card kind, so no caller needs to pre-warm the cache.

## How to use

- `GetCard(id)` — card for an ID; panics on the `Invalid` sentinel or out-of-range.
- `CardByName(displayName)` — `(CardID, false)` when not found. Keyed on `DisplayName`
  ("Aether Slash [R]"), so each pitch variant is a distinct entry.
- `AllCards()` — every valid card ID in registration order; freshly allocated, safe to mutate.
- `HeroByName(name)` / `WeaponByName(name)` — `(nil, false)` when not registered.
- `AllWeapons` — every implemented weapon, for loadout enumeration.
- `Registry{}.LegalCards()` / `.LegalWeapons()` — the filtered deck-construction pools.

## How to add an entry

- A new card is added automatically: `cmd/cardgen` regenerates `cards_gen.go` (the `cardsByID`
  slice) from the card YAML files, so adding a card and running `go generate` registers it.
- A new hero is added to the `all` slice in `heroesByName` (`heroes.go`).
- A new weapon is added to `AllWeapons` (`weapons.go`).

## Important files

- `registry.go` — package doc, the `Card` / `Hero` / `Weapon` and marker interfaces.
- `cards.go` — `CardID` alias, `cardsByName` index, `GetCard` / `CardByName` / `AllCards`.
- `cards_gen.go` — generated `cardsByID` slice; do not hand-edit.
- `heroes.go` — `heroesByName`, `HeroByName`.
- `weapons.go` — `AllWeapons`, `weaponsByName`, `WeaponByName`.
- `interfaces.go` — `Registry`, `LegalCards`, `LegalWeapons`, `isExcludedFromPool`.

## Gotchas and invariants

- `cards_gen.go` is generated — never hand-edit it; regenerate via `go generate
  ./internal/card/cards/...`.
- `cardsByID` index 0 is `nil` (the `Invalid` sentinel); iteration starts at 1 and must
  null-check.
- Name indexes are keyed on `DisplayName`, not bare `Name`, so the three printings of a card
  resolve to distinct IDs.
- The `NotSilverAgeLegal` set must stay in sync with `data_sources/silver_age_banlist.txt`.
- The registry and the sim never import each other; cross-package contracts are matched
  structurally.
