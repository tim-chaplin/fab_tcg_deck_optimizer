# Refactor design: game-vs-optimizer separation

## Status

Draft. Open questions throughout — please mark resolutions inline so we converge on
the layout before any code moves.

## Why

The codebase has accumulated cross-cutting modules. `TurnState` is a big bag holding
*both* game state (deck, hand, graveyard, runechant tokens, in-play auras) and
turn-scoped transients (value accumulator, log entries, incoming damage, pitched-this-turn).
Cards mutate every field of it directly. Hand-evaluator code lives in `internal/hand`
but reaches into `card.TurnState` for game state. Deck-evaluator code lives in
`internal/deck` but mutates the same TurnState fields.

This makes it hard to:
- Reason about which code is allowed to mutate which state.
- Add cross-cutting concerns like cacheable-tracking without footguns.
- Tell at a glance whether a given file is modelling FaB's rules or modelling our
  optimizer's loop.

## Goals

1. **Each game object gets a class**, responsible only for understanding and managing
   its own state.
2. **Sharp separation** between modules that *implement game concepts* (FaB rules) and
   modules that *implement our concepts* (the optimizer's eval loop, value scoring,
   annealing, statistics).
3. **Encapsulation that makes invariants enforceable** — once deck/graveyard live in
   real classes with method APIs, cacheable-tracking can hang naturally off the read
   accessors instead of needing static lints to police direct field access.

## Glossary: game vs optimizer

**Game concepts** (the FaB rules):
- Cards, decks, hands, heroes, weapons, in-play auras, tokens (Runechants), zones
  (graveyard, banish, arsenal).
- The act of playing a card: a card's `Play` mutates game state.
- Triggers fired when a game event happens (hero ability, aura "when X is played").

**Optimizer concepts** (our infrastructure, sit *outside* the game):
- The hand evaluator that searches all partitions / chain orderings for the highest-
  value line.
- The deck evaluator that runs N shuffles, accumulates per-turn value, and reports
  statistics.
- Annealing / iterate modes that mutate decks in search of higher mean value.
- Mutation enumerators (single-card swaps, pair swaps, weapon-loadout swaps).
- Caches (the future hand-eval cache that started this whole thread).
- Format-legality predicates, deck I/O.

## Proposed package layout

```
internal/
├── game/
│   ├── card/                 - Card interface, CardState, ID, TypeSet, markers,
│   │   │                       LikelyToHit, DisplayName, chain-step text cache
│   │   ├── runeblade/        - Runeblade card implementations (subpackage)
│   │   ├── generic/          - Generic card implementations
│   │   └── fake/             - Test fakes
│   ├── deck/                 - Deck struct: collection of cards, Shuffle / Draw /
│   │                           Peek / Pop / Prepend / Tutor / Random / Sideboard
│   ├── hand/                 - Hand struct: cards a player holds
│   ├── hero/                 - Hero interface + concrete (Viserai, …); ID type
│   ├── weapon/               - Weapon interface + concrete (NebulaBlade, …)
│   ├── zone/                 - Graveyard, Banish, Arsenal (each its own type)
│   ├── token/                - Runechant pool (or similar tokens later)
│   ├── aura/                 - AuraTrigger + in-play-aura concept
│   └── state/                - GameState aggregate + Turn transients
│
└── opt/
    ├── handeval/             - Hand evaluator (was internal/hand): partition
    │                           enumeration, sequence enumeration, chain replay
    ├── deckeval/             - Deck evaluator (was internal/deck/evaluate.go):
    │                           N-shuffle loop, per-turn fold, adaptive stop
    ├── anneal/               - Anneal loop, iterate loop, mutation enumerator
    ├── stats/                - Stats type, histogram, marginal-card stats
    ├── deckio/               - Deck JSON / .txt I/O (kept where it is, possibly)
    ├── deckformat/           - Format predicates
    └── fabrary/              - Fabrary import
```

`cmd/fabsim/` is unchanged conceptually — just imports rewire.

### Open question 1: depth of the `game/` and `opt/` hierarchy

Two flat groupings (`game/...`, `opt/...`) reads cleanly but adds one path segment to
every import. Alternative: keep packages at `internal/<name>/` (no `game/` / `opt/`
prefix) and just rely on naming. Vote?

### Open question 2: where does `card/<subdir>/` go?

The user's vision puts card implementations under `card/`. That nests subpackages
two levels deep under `game/` (`internal/game/card/runeblade/`). Acceptable, or
should card subpackages live elsewhere?

## Game-object classes — what each owns

### Card

**Owns:** name, ID, cost, pitch, attack, defense, types, go-again-printed flag,
Play behavior.

**Already has its own class** ✓ — the card subpackage types implement the `Card`
interface. The `Play` hook receives shared state and mutates it. Cards do *not*
manage cross-card state themselves.

**Open question 3:** the `Card` interface currently has `Play(s *TurnState, ...)`.
After the split, where does `Play` live? Three options on the table:

- (a) `Card.Play` keeps its current shape, but `TurnState` is renamed `GameState` and
  lives in `game/state/`. Cards still take the whole game state. Simple but cards
  retain access to everything.
- (b) `Card.Play` takes a narrower context (e.g. `Play(g *Game)` where Game is the
  per-game aggregate). Forces cards to go through the aggregate's API rather than
  reach into individual zones — but ergonomically Game.Hand().Append(c) is a bit
  noisier than s.Hand = append(s.Hand, c).
- (c) `Card.Play` takes scoped objects: `Play(hand Hand, deck Deck, gy Graveyard, …)`.
  Most explicit but the parameter list is long.

### Deck

**Owns:** the cards in the deck, top-to-bottom order.

**Operations** (game-rule operations only — no eval-loop concerns):
- `Shuffle(rng)`
- `Top()` — peek top card
- `Pop()` — remove and return top card
- `Prepend(c)` — push to top
- `Tutor(predicate)` — find and remove first match
- `Len()`
- `Cards()` — all cards (for serialisation / inspection)

The current `internal/deck` package mixes Deck-the-game-object with deck-evaluation
logic (Evaluate, IterateParallel). The eval logic moves out. What stays:
- `Deck` struct
- `New(hero, weapons, cards) *Deck`
- `Random(hero, size, maxCopies, rng, legal)` — random deck construction (game-side)
- `Sideboard`, `Equipment` fields (already game-side)
- `LegalPool` filtering

**Open question 4:** the current `Deck` struct also embeds `Stats` (the optimizer's
result type). That's a category violation — Stats is an optimizer concept. Move
`Stats` out of Deck (Stats lives separately, callers pair them when needed)?

### Hand

**Owns:** the cards a player holds during a turn.

**Operations:**
- `Add(c)` — append (drawn card, tutored card)
- `Remove(c)` — pop a specific card (alt-cost, pitch)
- `Cards()` — read-only view
- `Len()`

**Open question 5:** the current `TurnState.Hand` is a `[]card.Card` slice. Is the
right answer to extract it into a real Hand object, or is "hand is just a slice"
the simplest model and we don't need a class for it?

### Hero

**Owns:** name, ID, intelligence (handsize), health, types, OnCardPlayed hook.

**Already has its own class** ✓ — `Viserai` etc. implement the `Hero` interface. The
interface lives in `internal/hero/`. After the split, lives in `internal/game/hero/`
(or stays at `internal/hero/`).

The `OnCardPlayed` hook references TurnState. Refactor needs to reconcile that:

- (a) OnCardPlayed takes the new per-game state aggregate.
- (b) OnCardPlayed takes scoped pieces (the relevant zones).

### Weapon

**Owns:** name, ID, cost, attack, defense, durability(?), Play hook.

**Already has its own class** ✓ — Nebula Blade etc. Currently in `internal/weapon/`.
Same questions as Card around the Play hook's parameter shape.

### Zones (Graveyard, Banish, Arsenal)

**Currently** `[]card.Card` slices on TurnState — Graveyard/Banish — and a single
Card field — Arsenal.

**Proposed:**
- `Graveyard` — methods: `Add(c)`, `Cards()`, `BanishMatching(pred) (Card, bool)`,
  `Len()`. Read accessors are the cacheable-tracking flip points (scanning the
  graveyard reads prior-turn hidden state).
- `Banish` — methods: `Add(c)`, `Cards()`. Append-only typically.
- `Arsenal` — single-card slot; methods: `Get()`, `Set(c)`, `IsEmpty()`.

### Tokens

**Currently** a single `Runechants int` counter on TurnState plus an
`ArcaneDamageDealt` flag.

**Proposed:** `RunechantPool` (or just `Runechants`) type with `Create(n)`,
`ConsumeAll() int`, `Count()`, `ArcaneDamageDealt()`, `MarkArcaneDamageDealt()`.

**Open question 6:** is `Runechants` enough of a class to deserve its own type, or
is the int counter + bool flag the right level of abstraction?

### Auras

**Currently** `AuraTrigger` structs in `card.TurnState.AuraTriggers`. Auras-as-cards
(Sigil of Suffering, Sigil of the Arknight) are just Card implementations that, on
Play, register an AuraTrigger.

**Proposed:** `AuraSet` type holding `[]AuraTrigger` with methods `Add(t)`,
`FireMatching(type, ...)`, `Count(type)`, etc. The framework's
`fireAttackActionTriggers` and start-of-turn fire loop become methods on AuraSet.

### Turn (transient state)

The remainder of `TurnState` — Value, Log, CardsPlayed, AuraCreated, IncomingDamage,
Pitched, EphemeralAttackTriggers, Revealed, TriggeringCard, NonAttackActionPlayed,
SkipLog, Overpower — is per-chain / per-turn transient. These are the
"this-chain-step" bookkeeping bits.

**Proposed:** `Turn` type aggregating these. Methods: `RecordValue(n)`,
`AddLogEntry(...)`, `RecordPlayed(c)`, `BlockIncoming(n)`, `RegisterEphemeral(t)`,
etc. (most of these methods exist on TurnState today).

### GameState aggregate

The composite that cards' Play hooks receive. Holds (or accesses) all of the above:

```go
// rough sketch
type GameState struct {
    Hero      Hero
    Deck      *Deck
    Hand      *Hand        // (or just []Card if open-question-5 says no)
    Graveyard *Graveyard
    Banish    *Banish
    Arsenal   *Arsenal
    Tokens    *RunechantPool
    Auras     *AuraSet
    Turn      *Turn
}
```

`Card.Play(g *GameState, self *CardState)`. Cards go through `g.Deck.Pop()`,
`g.Graveyard.Add(c)`, `g.Tokens.Create(1)`, `g.Turn.AddLogEntry(...)`, etc.

## Optimizer-side modules — what each owns

### handeval

**Owns:** the partition / sequence / chain-replay search that, given a hand against
incoming damage and a deck snapshot, returns the highest-value `TurnSummary`.

**Inputs:** Hero, Weapons, Hand cards, IncomingDamage, Deck snapshot, runechant
carryover, arsenal-in card, prior aura triggers.

**Outputs:** `TurnSummary` (Value, BestLine, end-of-chain CarryState).

Was `internal/hand/`. Moves to `internal/opt/handeval/` (or wherever the optimizer
tree lives).

### deckeval

**Owns:** the N-shuffle loop. For each shuffle, deal hands and call handeval per turn,
fold per-turn results into running statistics, until the deck runs out or the early-
stop policy fires.

**Inputs:** Deck (the *game object*), incomingDamage, RNG, eval-mode (fixed shuffles
vs adaptive).

**Outputs:** `Stats` (mean value, per-cycle averages, best-turn snapshot, marginal-
card stats, …).

Was `internal/deck/evaluate.go` + helpers. Moves to `internal/opt/deckeval/`.

### anneal / iterate

**Owns:** the simulated-annealing loop and the iterate loop. Both mutate decks in
search of higher mean value via deckeval.

Was `internal/deck/iterate.go` + the cmd/fabsim/mode_anneal.go logic.

### Mutations

**Owns:** the enumerator that lists candidate deck mutations (single-card swaps,
pair swaps, weapon-loadout swaps).

Was `internal/deck/mutations.go`, `weapon_loadouts.go`, `card_pairs.go`.

**Open question 7:** mutations operate on Decks but are an optimizer concept (they
exist for the search loop). Goes under `opt/mutations/`?

### Stats

**Owns:** the `Stats` type, the value histogram, per-card marginal-stats accounting,
mean / variance / SE computation.

Was `internal/deck/stats.go`. Moves to `internal/opt/stats/`.

## Migration approach

Doing this in one merge is too risky and too much code-review at once. Proposal:

1. **PR-1: Lift TurnState into per-zone types** — within `internal/card`, replace
   the bag-of-fields TurnState with a struct that has typed sub-fields (`Deck`,
   `Hand`, `Graveyard`, etc.). Cards still get the same TurnState parameter; the
   surface change is `s.Deck = ...` → `s.Deck.Pop()` etc. No package moves yet.
   Lets us prove the per-zone API without paying the import cost.
2. **PR-2: Move handeval out** — `internal/hand` → `internal/opt/handeval`
   (or whatever final path). Pure rename + import update.
3. **PR-3: Move deckeval / anneal / mutations / stats out** — extract from
   `internal/deck` into `internal/opt/...`. `internal/deck` ends up holding only
   the Deck game-object.
4. **PR-4: Move Hero interface and concrete heroes** under `internal/game/hero` (if
   we agree on the `game/` prefix).
5. **PR-5: Move card-effect packages and remaining game pieces** under
   `internal/game/card/...` etc. (if we agree on the prefix).
6. **PR-6: Cacheable tracking** — now that deck/graveyard are encapsulated and
   `Play` goes through a typed API, add the IsCacheable bit and the read-flip
   semantics. (This is the original PR #230 work, redone on the cleaner foundation.)

### Open question 8: PR scope

Are these the right PR boundaries? In particular:
- Should PR-1 split per-zone *and* swap to a "scoped" Play signature
  (`Play(g *GameState, ...)` vs `Play(s *TurnState, ...)`)? Or two separate PRs?
- Does PR-2 / PR-3 really need to be two PRs, or merge into one "extract opt
  loop"?
- Where does the `game/` / `opt/` hierarchy decision sit (PR-1 or its own pre-PR)?

## Decisions to lock before code moves

Numbered list of the open questions above so we have a checklist:

- [ ] **Q1**: `game/` and `opt/` prefixes, or flat layout under `internal/`?
- [ ] **Q2**: where do card-effect subpackages (`runeblade`, `generic`, `fake`) live?
- [ ] **Q3**: what does `Card.Play`'s signature look like after the split?
- [ ] **Q4**: does `Stats` come off `Deck`?
- [ ] **Q5**: is `Hand` worth a class or is `[]Card` fine?
- [ ] **Q6**: is `RunechantPool` worth a class or is `int` fine?
- [ ] **Q7**: mutations live under `opt/`?
- [ ] **Q8**: are the proposed PR boundaries right?
