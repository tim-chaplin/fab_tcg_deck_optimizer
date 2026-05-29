# internal/sim

## Purpose

`sim` is the hand-and-deck evaluator at the heart of the optimizer. Given a 40-card deck it
shuffles, walks two cycles of hands, and for each hand brute-forces the optimal turn line:
the partition of the hand into roles (Pitch, Attack, Defend, Held, Arsenal) and the ordering
of the attack-turn search that maximises turn value. Turn value is damage dealt plus damage
prevented. Per-hand results fold into `deck.Stats`, and `RunMutationRound` applies a
simulated-annealing acceptance gate so a search loop can climb toward high-value decks.

## The two-layer search

`(*Evaluator).findBest` (in `partition.go`) runs the search as two nested layers:

1. **Partition enumeration** (`partition.go`) — a recursive walk assigns every hand card (and
   the optional arsenal-in card, treated as index `n`) one of the five roles. `roleAllowed`
   prunes illegal assignments; the `Defend` role is skipped entirely when `IncomingDamage`
   is 0. Each complete leaf is handed to `evaluatePartition` → `bestAttackWithWeapons`.
2. **Attack-turn search** (`sequence.go`) — for one partition leaf, `bestAttackWithWeapons`
   enumerates the phase-split mask (`pmask`, splitting pitch resources across the attack and
   defense phases when defense reactions are present) and the weapon/item-ability mask
   (`wmask`, choosing which weapons and item abilities to swing). For each mask combination
   `bestSequence` permutes the attacker list via Heap's algorithm (and the attack-pitch
   ordering, and the cartesian product of modal modes), replaying each ordering through
   `playSequenceWithMeta`.

`playSequenceWithMeta` is the attack-turn runner: it builds a fresh per-permutation `*GameState`,
removes each played card from hand, funds it through the `pitchPool`, spends an action point
unless the step is free, fires `CardOrAbility` triggers, calls `ResolveAttackStep`, finalizes
each attack's on-hit effects, and fires `EndOfTurn`. Every leaf and permutation is scored by
`attackTurnScore` — a lexicographic tuple of `(value, cardsPlayed, totalCards, totalCounters)` —
and the highest-ranked is kept.

## Entry points

- **`best` / `(*Evaluator).Best`** (`hand_eval.go`) — package-private and method forms of
  the single-turn optimiser. `Best` checks the eval cache, then runs `findBest`. Production
  callers route through the test entry points or `Evaluate`.
- **`(*Evaluator).Evaluate`** (`deck_eval.go`) — shuffles a deck `runs` times, walks two
  cycles of hands per shuffle, threads cross-turn carryover, and returns `deck.Stats`.
- **`EvalOneTurnForTesting` / `EvalTwoTurnsForTesting`** (`integration_testing.go`) — the
  public single-/two-turn test entry points. They drive the same per-turn pipeline production
  uses (no shuffle) and return the resulting `TurnSummary`. Turn-level tests in `turntests/`
  use these and `(*Evaluator).Evaluate` exclusively.
- **`RunMutationRound`** (`mutation_round.go`) — evaluates a batch of candidate decks under
  the Metropolis acceptance rule and returns the first accepted mutation.

## Key types

- **`Evaluator`** — caches per-goroutine scratch (`attackBufs`) and an optional hand-eval
  cache. Not safe for concurrent use; each goroutine constructs its own. Constructors:
  `NewEvaluator`, `NewEvaluatorParallel`, `NewEvaluatorWithCache`, `NewEvaluatorWithoutCache`.
- **`TurnSummary`** (`hand_types.go`) — the result of one turn: the winning `BestLine` of
  role assignments, swung weapon names, `Value`, and the post-attack-turn `*GameState`.
- **`Matchup`** (`matchup.go`) — the opponent-profile constants (`IncomingDamage`,
  `ArcaneIncomingDamage`) held constant for an `Evaluator`'s lifetime.
- **`Cache` / `CacheStats`** (`eval_cache.go`) — the thread-safe hand-eval cache and its
  counter snapshot.

## Per-goroutine scratch buffers

`attackBufs` (`attackbufs.go`) pools every slice and pooled `*GameState` / `*GameEngine` /
`*Deck` the search reuses, sized once per `(handSize, weapons)` shape and amortised across
every hand a long-running iterate pass evaluates. `cardmeta.go` holds `attackerMeta` — a
lazily-populated, ID-keyed table of scalar card attributes (types, cost bounds, GoAgain,
attack/DR/modal flags) so the attack turn inner loop avoids interface dispatch. The metadata table
is the one piece of shared state safe to read from all goroutines (written once per ID under
a mutex, then read lock-free).

## The eval cache

The hand-eval cache (`eval_cache.go`) memoizes the winning solution per `evalCacheKey` — the
sorted hand multiset plus loadout, carryover auras/items, and arsenal/marked state. On a hit,
`replayBest` (`eval_cache_replay.go`) rebuilds the `TurnSummary` by replaying the captured
attacker order, modal modes, pitch order, and blocker modes verbatim — no partition search,
no permutation enumeration. Storing is gated on `Cacheable`: if any sibling partition touched
hidden state (deck or graveyard via an engine accessor), the result depends on state the key
doesn't capture and is not stored. `Matchup` is intentionally absent from the key — an
`Evaluator`'s lifetime spans a constant matchup; tests mixing matchups must use
`NewEvaluatorWithoutCache`.

## ResolveAttackStep

`ResolveAttackStep` (on `*gameengine.GameEngine`, called by the attack-turn runner and the defense
pass) runs `card.Play`, then credits `g.Value` and appends the canonical
`<Card>: <VERB> (+N)` attack-step log entry. `N` is `EffectiveAttack` for attacks and weapon
swings, the capped `EffectiveDefense` (with the `IncomingDamage` decrement) for defense
reactions and `DefensiveInstant` cards, and 0 otherwise. Because the log entry is appended
*after* `Play` returns, any self-buff `Play` applied is reflected in the displayed delta.

## Logging idioms

A card's `Play` body does **not** emit its own attack step — `ResolveAttackStep` owns that.
`Play` emits **rider sub-lines only**: post-triggers attributed to `self.Card.DisplayName()`
for self-riders, or to a different source for cross-card riders (e.g. an `OnHit` attached to a
target card). The framework threads a `card.Logger` into every Card-shaped hook; a nil-pointer
logger (the find-best skip pass) silently elides every call, so cards never gate logging
manually.

**A line starting with `l.AppendXxx(...)` must have no side effects.** Put the state change
on its own preceding line so a reader scanning for "what does this card do" can skip every
`l.AppendXxx` line:

```go
// Good — the Append line is pure:
s.CreateRunechants(2)
l.AppendPostTrigger(self.Card.DisplayName(), "Created 2 runechants", 2)

// Bad — side effect hidden inside the log call:
l.AppendPostTrigger(self.Card.DisplayName(), "Created 2 runechants", s.CreateRunechants(2))
```

Conditional self-buffs go in the `Play` body, *before* the attack step is logged —
`ResolveAttackStep` reads the buffed `EffectiveAttack` / `EffectiveDefense` after `Play`
returns:

```go
// Bluster Buff (mode 1): +2{p} for paying {r}.
func blusterBuffPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
    if self.Mode == 1 {
        self.BonusAttack += 2
    }
}
```

## Important files

- `partition.go` — `findBest`, the partition recursion, `roleAllowed`, `defendersDamage`,
  `attackTurnScore`.
- `sequence.go` — `bestAttackWithWeapons`, `bestSequence`, `playSequenceWithMeta`, the
  per-permutation state lifecycle.
- `hand_eval.go` — `Best`, the `Evaluator` type and its constructors.
- `deck_eval.go` — `Evaluate`, the shuffle loop, cross-turn carryover, per-turn stats.
- `integration_testing.go` — `EvalOneTurnForTesting` / `EvalTwoTurnsForTesting`.
- `attackbufs.go` — the pooled scratch struct and `fillPartitionCards`.
- `cardmeta.go` — the per-card metadata cache and `costAt`.
- `eval_cache.go` / `eval_cache_replay.go` — the hand-eval cache and verbatim replay.
- `evaluate_partition.go` — the partition-leaf → attack-turn-search bridge.
- `mutation_round.go` — `RunMutationRound` and the simulated-annealing acceptance gate.
- `matchup.go`, `hand_types.go` — `Matchup` and `TurnSummary`.
- `print_best.go`, `format.go`, `turnlog.go` — display rendering for the winning line.

## Gotchas and invariants

- Hands are sorted by `Card.ID()` before search so cache-on and cache-off paths produce
  byte-identical results for matching multisets.
- Each attack-turn permutation runs against a fresh per-permutation `*GameState` borrowed from
  the Evaluator's prewarmed pool. Each pool slot owns its own hand / cardsPlayed /
  graveyard / banished / deck-wrapper backings, so the winning slot's state survives
  unmolested when the next permutation borrows a different slot.
- The attack-budget prune relaxes by `maxResourceBonus` (the declared upper bound of
  resource-producing cards) so resource cards aren't pruned out; `pitchPool.pay` does the
  exact funding check.
- Cards with state-dependent cost must implement `VariableCost.EffectiveCost(g)` so the
  attack-turn runner reads the live figure at `costAt`. The interface contract on
  `Card.Cost()` (no engine arg) makes the legacy "secretly-varying static cost" footgun
  impossible by construction.
- New `turntests/` files must not call `ResolveAttackStep` directly — they exercise the
  public `Eval*` entry points.
</content>
