# fab-deck-optimizer

A deck-building and simulation tool for the Flesh and Blood TCG, written in Go.

Built with Claude Opus 4.6.

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

`fabsim` has four run modes, selected with `-mode`:

All modes read and write `mydecks/<deck>.json` where `<deck>` comes from `-deck` (default
`best_deck`). The `.json` suffix on `-deck` is optional.

- **`random`** (default) — two-phase search. Generates `-decks` random decks and evaluates each
  shallowly (`-shallow-shuffles` shuffles); takes the top `-top-n` and re-evaluates them with more
  shuffles (`-deep-shuffles`). Writes the winner to the deck file if it beats whatever's already
  there.
- **`iterate`** — loads the deck file and hill-climbs on it deterministically: each round
  enumerates every single-slot mutation (every alternative weapon loadout + every (card-in-deck,
  card-out-of-deck) swap), adopts the first one that scores higher, and restarts. When a full
  round finishes without finding an improvement, the deck is at a local maximum and `iterate`
  exits. Press Enter to abort mid-round. If the deck file doesn't exist yet, `iterate` bootstraps
  with a single random deck and climbs from there — you don't have to run `random` first.
- **`eval`** — loads the deck file, simulates it for `-deep-shuffles` hands against `-incoming`
  damage, and prints the resulting stats. Does **not** overwrite the file — use this to re-score a
  saved deck at a new shuffle depth or opponent pressure without clobbering whatever's on disk.
- **`print`** — prints the deck without running any simulation.

### Suggested workflow

Start with a wide random search to seed `mydecks/best_deck.json`, then hill-climb from there:

```
go run ./cmd/fabsim -mode=random  -decks=10000 -shallow-shuffles=10 -top-n=100 -deep-shuffles=1000
go run ./cmd/fabsim -mode=iterate -deep-shuffles=1000
```

`random` explores the space; `iterate` refines the best find. Re-run either stage as often as you
like — each run only overwrites `mydecks/best_deck.json` if it finds something better.

### Flags

- `-mode` — `random`, `iterate`, `eval`, or `print` (default `random`)
- `-decks` — number of random decks to generate in phase 1 of `random` (default 10000)
- `-shallow-shuffles` — shuffles per deck in phase 1 wide search (default 10)
- `-top-n` — number of phase-1 decks to advance to phase 2 (default 100)
- `-deep-shuffles` — shuffles per deck in phase 2 deep eval (also used per mutation in `iterate`)
  (default 1000)
- `-incoming` — opponent damage per turn (default 4)
- `-deck-size` — cards per deck (default 40)
- `-max-copies` — max copies of any single card printing (default 2)
- `-seed` — RNG seed (default: time-based)
- `-deck` — deck name; resolved to `mydecks/<name>.json` (default `best_deck`). The `mydecks/`
  directory is created automatically.

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
