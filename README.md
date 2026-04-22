# fab-deck-optimizer

A deck-building and simulation tool for the Flesh and Blood TCG, written in Go.

Built with Claude Opus 4.6.

## FAQ

### So are these AI-generated decks?

No. This is just a computer program that implements an evaluation function and uses known
optimization techniques to look for optimal decks according to the evaluation. You can compile
and run the program to find decks on your own computer without having invoked AI at all.

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

- Loads a 40-card deck from `mydecks/<name>.json`.
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
`<hero>_<format>_<incoming>_incoming`, e.g. `viserai_silver_age_0_incoming`, so different (hero,
format, `-incoming`) regimes keep separate deck files). The `.json` suffix on `-deck` is
optional.

- **`anneal`** — simulated-annealing search on the deck at `-deck`, or on a fresh random deck
  if the file doesn't exist yet. Each round enumerates every single-slot mutation (every
  alternative weapon loadout + every (card-in-deck, card-out-of-deck) swap). Mutations are
  screened at `-shallow-shuffles`; candidates that clear the acceptance gate are re-evaluated
  at `-deep-shuffles` to confirm. The acceptance gate is the Metropolis rule: strict
  improvements are always accepted, worse mutations are accepted with probability
  `exp((avg - baseline) / T)` when `-start-temp > 0`. Temperature decays geometrically per
  acceptance (`-temp-decay`, floored at `-min-temp`). At `-start-temp 0` the gate is strictly
  `> baseline` and anneal degenerates to a classical hill climb. A round with zero acceptances
  is treated as a local maximum and anneal exits. Press Enter to abort mid-round (exits 130
  so wrapper scripts can tell this apart from natural convergence). Only the best-ever deck is
  persisted to disk — walks through worse states under annealing don't regress the JSON.
- **`eval`** — loads the deck file, simulates it for `-deep-shuffles` hands against `-incoming`
  damage, and prints the resulting stats. Does **not** overwrite the file — use this to re-score a
  saved deck at a new shuffle depth or opponent pressure without clobbering whatever's on disk.
- **`print`** — prints the deck without running any simulation.
- **`import`** — interactively imports a deck from fabrary.net. Prompts for a deck name, then
  asks you to paste the plain-text export; input ends automatically at fabrary's
  `See the full deck @ …` footer. Saves the result as `mydecks/<name>.json`. The `-deck` flag is
  ignored — the name always comes from the prompt. Cards the optimizer hasn't implemented yet are
  skipped with a warning rather than blocking the import.
- **`diff`** — prints the card-count delta between two saved decks. Usage: `fabsim diff <deck1> <deck2>`.

### Suggested workflow

Start with annealing to escape weak local maxima, then re-anneal repeatedly to probe the
neighbourhood of each new best:

```
go run ./cmd/fabsim anneal -start-temp 1 -incoming 7
./scripts/anneal-reanneal.ps1 -Deck viserai_silver_age_7_incoming -StartTemp 1 -Incoming 7
```

Or fan out across several independent starts and rank the results:

```
./scripts/anneal-restarts.ps1 -N 10 -DeckTemplate 'viserai_*' -Incoming 7 -StartTemp 1
```

Each run only overwrites its deck file when a new best-ever avg is found, so it's safe to
re-run indefinitely.

### Flags

- `-shallow-shuffles` — shuffles per deck when screening anneal mutations (default 100)
- `-deep-shuffles` — shuffles per deck when confirming anneal improvements and for `eval` (default 10000)
- `-incoming` — opponent damage per turn (default 0)
- `-deck-size` — cards per deck (default 40)
- `-max-copies` — max copies of any single card printing (default 2)
- `-seed` — RNG seed (default: time-based)
- `-deck` — deck name; resolved to `mydecks/<name>.json` (default
  `<hero>_<format>_<incoming>_incoming`, keyed off the hero, format, and `-incoming`). The
  `mydecks/` directory is created automatically.
- `-format` — constructed format whose banlist restricts the card pool during search. Defaults
  to `silver_age`, which is currently the only supported format. The authoritative Silver Age
  banlist lives at `data_sources/silver_age_banlist.txt`.
- `-start-temp` — anneal: starting temperature. `0` (default) runs a pure hill climb. Higher
  values probabilistically accept worse mutations early (Metropolis rule).
- `-temp-decay` — anneal: multiplicative cooling per acceptance (default 0.95).
- `-min-temp` — anneal: temperature floor (default 0).
- `-finalize` — anneal: high-precision pass — overrides `-shallow-shuffles` to 10000 and
  `-deep-shuffles` to 100000. Use on a deck that's already converged to squeeze out the
  remaining sub-percent improvements.
- `-reevaluate` — anneal: force re-evaluation of the loaded deck's baseline avg even if its
  prior run count already matches `-deep-shuffles`. Use after adjusting modelling assumptions.

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
.github/workflows/   GitHub Actions (anneal-sweep: parallel matrix run)
cmd/fabsim/          CLI entry point
cmd/parsecarddb/     Card-database parser / filter
internal/card/       Card interface, TurnState, and card implementations
internal/deck/       Deck construction and parallel anneal round driver
internal/deckio/     JSON serialisation / deserialisation
internal/hand/       Optimal-play solver for a single hand
internal/hero/       Hero definitions and on-play triggers
internal/weapon/     Weapon definitions
scripts/             PowerShell wrappers for multi-restart and reanneal sweeps
```
