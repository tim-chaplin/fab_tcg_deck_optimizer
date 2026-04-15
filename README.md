# fab-deck-optimizer

A deck-building and simulation tool for the Flesh and Blood TCG, written in Go.

## What it does today

A minimal first cut:

- A 40-card deck (currently hardcoded: 20 generic blue + 20 generic red).
- Simulates shuffling the deck and repeatedly drawing hands of 4 cards.
- For each hand, solves for the optimal play in isolation: partition the hand
  into Pitch / Attack / Defend roles, where pitch resources must cover
  attacker costs, and hand value = damage dealt + damage prevented (capped at
  the opponent's incoming damage per turn).
- Per FaB rules, pitched cards return to the bottom of the deck; attacked
  and defended cards are spent. The simulation runs until fewer than 4
  cards remain.
- Reports the overall average hand value, plus the averages for the first
  and second cycle through the deck.

## Cards

The toy card set:

| Card | Cost | Pitch | Attack | Defend |
|------|-----:|------:|-------:|-------:|
| Blue | 1    | 3     | 1      | 3      |
| Red  | 1    | 1     | 3      | 1      |

## Usage

```
go run ./cmd/fabsim -runs=10000 -incoming=4
```

Flags:

- `-runs` — number of shuffles (default 10000)
- `-incoming` — opponent damage per turn (default 4)
- `-seed` — RNG seed (default: time-based)

## Tests

```
go test ./...
```

## Layout

```
cmd/fabsim/       CLI entry point
internal/card/    Card type and the two example cards
internal/hand/    Optimal-play solver for a single hand
internal/sim/     Deck simulation and stat aggregation
```
