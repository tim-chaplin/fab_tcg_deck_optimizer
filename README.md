# fab-deck-optimizer

A deck-building and simulation tool for the Flesh and Blood TCG, written in Go.

## Goal

Find optimal deck lists for **goldfishing** — i.e. maximizing a deck's own output in a vacuum,
without modeling a live opponent. The simulator partitions each drawn hand into its best Pitch /
Attack / Defend split and reports aggregate value across many shuffles.

## Scope & limitations

This is a work in progress. The current model is deliberately narrow:

- **Hero pool.** Only cards legal for **Viserai in the Silver Age** are in scope (plus Generic
  cards). Other heroes / talents / formats aren't modeled yet.
- **Turns are evaluated in isolation.** There is no between-turn state — no arsenal, no persistent
  auras carrying over, no health totals, no deck thinning effects that span turns. Each hand is
  solved as a standalone puzzle.
- **No opponent counterplay.** The opponent is represented by a single configurable `-incoming`
  value: a static amount of damage per turn that the hand can defend against. There are no blocks
  from hand, no disruption, no reaction windows. This is the goldfishing assumption.
- **Simplified card effects.** Conditional bonuses are modeled where tractable (e.g. Runechant
  tokens count as +1 damage, "if hits" is assumed, "next Runeblade attack" riders peek forward via
  `CardsRemaining`). Effects that require deck / graveyard / multi-turn state are approximated as 0
  or omitted.
- **Card coverage is incomplete.** Most of the Runeblade Silver-Age pool and some Generics are
  implemented; the rest are stubbed or missing.

## What it does today

- Loads a 40-card deck (currently hardcoded).
- Shuffles and repeatedly draws hands of 4 cards.
- For each hand, brute-forces the optimal play: every partition of the hand into Pitch / Attack /
  Defend, every weapon-swing subset, every legal attack ordering (respecting Go again). Hand value
  = damage dealt + damage prevented (capped at `-incoming`).
- Per FaB rules, pitched cards return to the bottom of the deck; attacked and defended cards are
  spent. The simulation runs until fewer than 4 cards remain.
- Reports the overall average hand value, plus the averages for the first and second cycle through
  the deck.

## Usage

```
go run ./cmd/fabsim -shuffles=10000 -incoming=4
```

Flags:

- `-shuffles` — number of shuffles (default 10000)
- `-incoming` — opponent damage per turn (default 4)
- `-seed` — RNG seed (default: time-based)

Helper tool for exploring the upstream card database:

```
go run ./cmd/parsecarddb --names_only
```

## Tests

```
go test ./...
```

## Layout

```
cmd/fabsim/         CLI entry point
cmd/parsecarddb/    Card-database parser / filter
internal/card/      Card interface, TurnState, and card implementations
internal/deck/      Deck construction
internal/hand/      Optimal-play solver for a single hand
internal/hero/      Hero definitions and on-play triggers
internal/sim/       Deck simulation and stat aggregation
internal/weapon/    Weapon definitions
```
