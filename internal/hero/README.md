# internal/hero

## Purpose

Defines the `Hero` interface a Flesh and Blood hero card satisfies, plus the narrow engine
surfaces a hero ability needs. Concrete heroes live in the `internal/hero/heroes`
subpackage.

## Key types

- `Hero` — `ID`, `Name`, `Class`, `Types`, `Intelligence`, `OnCardPlayed`, and `Opt`. The
  view of a hero callers need beyond what the engine itself uses.
- `GameEngine` / `Logger` (interfaces.go) — the narrow engine and log surfaces hero
  abilities consume. They are declared package-locally so `hero` doesn't depend on
  `internal/card.GameEngine`; `*gameengine.GameEngine` satisfies them structurally. A hero
  that needs a richer surface type-asserts past the narrow interface (e.g. to
  `card.GameEngine` for class folding).

## How to use / extend

To add a hero, drop a `<hero>.go` file in `internal/hero/heroes/` declaring a struct that
implements `Hero`:

- `OnCardPlayed(played, ge, l)` runs whenever a card is played and returns the value the
  hero ability contributed. Viserai's implementation creates a Runechant when a Runeblade
  card follows a non-attack action this turn.
- `Opt(cards)` splits a revealed top-of-deck slice into `(top, bottom)` — the hero's
  card-ordering heuristic. Viserai's keeps one card per "slot category" (non-attack
  enabler, non-go-again action, block-only defender, blue pitch) and bottoms anything that
  over-fills an already-covered slot, so a balanced hand keeps feeding the Runechant
  trigger.

If a new hero ability needs more from the engine, extend the `GameEngine` interface in
`interfaces.go` rather than importing `internal/card` directly.

## Important files

- `hero.go` — the `Hero` interface and package doc.
- `interfaces.go` — the narrow `GameEngine` / `Logger` surfaces.
- `heroes/viserai.go` — Young Viserai: the only hero implemented today, the model for new
  hero files.

## Gotchas

- `OnCardPlayed` returns 0 for weapon swings: equipping or swinging a weapon isn't "playing
  a card" and doesn't trigger a hero like Viserai. The Viserai implementation checks for
  `TypeWeapon` explicitly.
- The `Opt` heuristic reads printed values only — `Types(nil)` and `GoAgain(nil)` skip any
  runtime-state fold, which is the right thing for an Opt-time decision.
