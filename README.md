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

`fabsim` takes a subcommand as its first argument. Running `fabsim` with no subcommand prints the
catalogue.

All subcommands read and write `mydecks/<deck>.json` where `<deck>` comes from `-deck` (default
`<hero>_<incoming>_incoming`, e.g. `viserai_0_incoming`, so different (hero, `-incoming`)
regimes keep separate deck files). The `.json` suffix on `-deck` is optional.

- **`random`** — two-phase search. Generates `-decks` random decks and evaluates each shallowly
  (`-shallow-shuffles` shuffles); takes the top `-top-n` and re-evaluates them with more shuffles
  (`-deep-shuffles`). Writes the winner to the deck file if it beats whatever's already there.
- **`iterate`** — hill-climbs deterministically on the deck at `-deck`, or on a fresh random
  deck if the file doesn't exist yet (so you don't need to run `random` first). Each round
  enumerates every single-slot mutation (every alternative weapon loadout + every (card-in-deck,
  card-out-of-deck) swap). Mutations are screened at `-shallow-shuffles`; only those that beat
  the current best on the shallow sample are re-evaluated at `-deep-shuffles` to confirm. The
  first mutation that beats the baseline at deep depth is adopted and the round restarts. When a
  full round finishes without finding a confirmed improvement, the deck is at a local maximum and
  `iterate` exits. Press Enter to abort mid-round.
- **`eval`** — loads the deck file, simulates it for `-deep-shuffles` hands against `-incoming`
  damage, and prints the resulting stats. Does **not** overwrite the file — use this to re-score a
  saved deck at a new shuffle depth or opponent pressure without clobbering whatever's on disk.
- **`print`** — prints the deck without running any simulation.
- **`import`** — interactively imports a deck from fabrary.net. Prompts for a deck name, then
  asks you to paste the plain-text export; input ends automatically at fabrary's
  `See the full deck @ …` footer. Saves the result as `mydecks/<name>.json`. The `-deck` flag is
  ignored — the name always comes from the prompt. Cards the optimizer hasn't implemented yet are
  skipped with a warning rather than blocking the import.

### Suggested workflow

Start with a wide random search to seed the deck file for the current `-incoming` setting, then
hill-climb from there:

```
go run ./cmd/fabsim random
go run ./cmd/fabsim iterate
```

`random` explores the space; `iterate` refines the best find. Re-run either stage as often as
you like — each run only overwrites the deck file if it finds something better.

### Flags

- `-decks` — number of random decks to generate in phase 1 of `random` (default 1000)
- `-shallow-shuffles` — shuffles per deck in `random` phase 1 and for screening `iterate`
  mutations (default 100)
- `-top-n` — number of phase-1 decks to advance to phase 2 (default 100)
- `-deep-shuffles` — shuffles per deck in `random` phase 2 and for confirming `iterate`
  improvements (default 10000)
- `-incoming` — opponent damage per turn (default 0)
- `-deck-size` — cards per deck (default 40)
- `-max-copies` — max copies of any single card printing (default 2)
- `-seed` — RNG seed (default: time-based)
- `-deck` — deck name; resolved to `mydecks/<name>.json` (default `<hero>_<incoming>_incoming`,
  e.g. `viserai_0_incoming`, keyed off the hero and `-incoming`). The `mydecks/` directory is
  created automatically.

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
