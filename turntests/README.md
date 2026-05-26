# turntests

## Purpose

Turn-level tests for the simulator. Each test drives a single turn (occasionally two) of
the engine from a fixed hand and asserts the resulting per-turn `Value` — the damage dealt
plus damage prevented the optimizer scores a hand by. Turn tests are the project's primary
way to verify card behaviour: they exercise a card through the real attack-turn runner, partition
search, and play ordering, so they catch interaction bugs a `Play`-direct unit test would
miss.

## What a turn test is

A turn test builds a `deck.Deck` (hero, weapons, deck list) and a hand of `card.Card`
values, runs one of the public Eval entry points, and asserts on the returned
`sim.TurnSummary`:

```go
func TestDefensiveInstant_BrushOffRedAlone(t *testing.T) {
    d := deck.New(heroes.Viserai, nil, fillerDeck())
    hand := []card.Card{cards.BrushOffRed{}}
    summary := sim.EvalOneTurnForTesting(d,
        gameengine.GameStateBuilder().SetIncomingDamage(5).Build(), hand)
    if got := summary.Value; got != 3 {
        t.Fatalf("Value = %d, want 3 (Brush Off Red prevents 3 of 5)", got)
    }
}
```

## How to add a turn test

1. Create `turntests/<card>_test.go` in `package turntests`.
2. Build a `deck.Deck` with a real hero from `internal/hero/heroes` (e.g.
   `heroes.Viserai`) and a hand of the cards under test, mixing in `internal/testutils`
   fakes for filler / known-value cards.
3. Drive the turn through a public Eval entry point only:
   - `sim.EvalOneTurnForTesting(deck, gameState, hand)` for one turn.
   - `sim.EvalTwoTurnsForTesting(...)` for cross-turn effects (auras, drawn-card payoffs).
   - `(*Evaluator).Evaluate` for a full multi-turn run.
   Build the `gameState` with the fluent `gameengine.GameStateBuilder()` (one `Set*` call
   per line when there are two or more).
4. Assert on `summary.Value`; use `sim.FormatBestLine(summary.BestLine)` in failure
   messages so a mismatch shows which line the engine actually picked.
5. Give the test a single brief doc sentence stating the behavior under test. Inputs and
   expected values are visible in the body.

## Public entry points only

`turntests/` is for public-entry-point tests. New files must **not** call
`ge.ResolveAttackStep(...)` directly — `internal/lint/turntests_lint_test.go` enforces this
via an allowlist of grandfathered files from the v2 migration. A test that needs an
unexported helper belongs as a same-package unit test next to the code instead, not here.
To migrate a grandfathered file, rewrite it against the public Eval API (or move it to a
same-package unit test) and remove its allowlist entry.

## Important files

- `doc.go` — the package doc.
- `mauvrion_skies_test.go`, `defensive_instant_test.go` — examples of the
  `EvalOneTurnForTesting` pattern, the model for new tests.

## Gotchas

- Every turn test must pass; a failing turn test is never traded away for speed — a
  fast-but-wrong model is worthless.
- Test a card's behaviour, not a transcription of its stat line: don't pin `GoAgain` /
  `Attack` / `Cost` for a specific card just to mirror the yaml.
