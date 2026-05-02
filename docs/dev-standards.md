# Developer Standards

Shared conventions for code and comments. Rules cited across multiple files factor in here so
the per-file comments can stay card-specific.

## Card file layout

Each card implementation lives in `internal/card/generic/` or `internal/card/runeblade/` (plus
`internal/card/fake/` for test doubles). A card file typically:

1. Declares a shared `TypeSet` var for the card's type line.
2. Declares one struct per printed pitch variant (e.g. `FooRed`, `FooYellow`, `FooBlue`).
3. Implements the `card.Card` interface plus any optional markers (`AddsFutureValue`,
   `VariableCost`, …) on each variant.
4. Shares a `fooPlay(...)` helper when variants differ only by a numeric parameter.

Card data (name, cost, pitch, attack, defense, type line, printed text) is transcribed from
`github.com/the-fab-cube/flesh-and-blood-cards`. Cards do not need a per-file `Source:` line;
the upstream repo is the authority for every card file in the project.

## Comment rules

- Wrap at 100 chars.
- Describe CURRENT behavior and its motivation. No history references ("replaces X", "was Y
  before", "now does Z"), no "previously/formerly/legacy/deprecated" framing.
- Delete completed TODOs instead of rewording them.
- Card docstrings cover card-SPECIFIC quirks — the printed rules text, any modelling fudge, and
  anything surprising about how this card interacts with the solver. Generic framework plumbing
  (how `AuraTrigger` is ticked, etc.) belongs in framework docs in `internal/card/card.go` and
  `internal/hand/hand.go`, not repeated in every card.
- Don't restate behavior that's already documented by an external function, type, or marker the
  card uses. If a card calls `card.LikelyToHit`, the docstring shouldn't re-explain the
  hit-likelihood heuristic; if a card carries `card.Dominator`, it shouldn't re-explain how
  Dominate interacts with `LikelyToHit`; if a card has a `NotImplemented` marker plus a
  `// not implemented: <quirk>` line, the docstring shouldn't repeat the same "rider isn't
  modelled" sentence in prose. Examples:
  - **Demolition Crew** (Generic Action - Attack with Dominate + an additional reveal cost) —
    no "Modelling: Dominate is advertised via the `card.Dominator` marker..." block. The
    `Dominator` interface implementation makes that link by itself; the additional reveal cost
    is documented by the `// not implemented:` comment above its `NotImplemented` method.
  - **Plunder Run** (a `// not implemented: on-hit draw rider...` line + a from-arsenal gate
    inside `Play`) — the docstring needs to call out the from-arsenal gate (card-specific
    quirk) but not the dropped on-hit draw (already on the marker).

## AuraTrigger lifecycle

Defined in `internal/card/card.go`. Standard shape for cards that "create an aura that fires
later":

- `Play` sets `s.AuraCreated = true` (so same-turn aura-readers see the aura) and calls
  `s.AddAuraTrigger(card.AuraTrigger{...})` with `Self`, `Type`, `Count`, and `Handler`.
- The sim walks `s.AuraTriggers` on each matching condition, invokes every matching `Handler`,
  decrements `Count`, and graveyards `Self` when `Count` hits zero.
- `OncePerTurn` caps an `AuraTrigger` at a single fire per turn.

Card docstrings should NOT restate this lifecycle. State only what's card-specific — the printed
clause, `Count = N`, and whatever the handler returns.

## NotImplemented vs Unplayable markers

Both markers exclude a card from random / mutation pools, so the optimizer skips them. They
mean different things:

- `sim.NotImplemented` — placeholder. The card *would* be worth modelling, we just haven't
  done it yet. Pair with a `// not implemented: <one-line description of the unmodelled
  rider>` comment immediately above the `NotImplemented()` method so the next implementation
  pass knows what's missing. Card files live in `internal/cards/notimplemented/`.
- `sim.Unplayable` — verdict. The card's effect is too weak to want even if fully modelled,
  so an implementation would be wasted work. The marker speaks for itself; **don't add a
  per-card rationale to the docstring**. Card files live in `internal/cards/unplayable/`.

The split keeps the unimplemented backlog honest: cards under `NotImplemented` are todos,
cards under `Unplayable` are closed. The directory split makes both lists visible at a
glance — `ls internal/cards/notimplemented/` is the live todo, and the lint test
`TestLayout_MarkersStayInSubpackages` enforces the layout so a stray marker can't silently
bypass the split.

## Standard rider wiring

Card docstrings should call out the printed rider and any modelling fudge, then stop. The
following plumbing is uniform and lives once in `internal/card/card.go`:

- **Played-from-arsenal go-again** (Fervent Forerunner, Frontline Scout, Performance Bonus,
  Promise of Plenty, Scour the Battlescape, …): cards call `self.GrantGoAgainIfFromArsenal()`
  at the top of `Play`; the helper flips `GrantedGoAgain` only when this copy came from the
  arsenal slot. Don't repeat the wiring per file — note that the rider only fires when this
  copy came from the arsenal slot.
- **+N{d} on arsenal-played defense reactions** (Springboard Somersault, Unmovable, …): cards
  implement `card.ArsenalDefenseBonus` and return `N`; `CardState.EffectiveDefense` folds the
  bonus in for the arsenal-in copy. Don't restate the wiring; just say "+N{d} when played from
  arsenal."
- **Conditional go-again / dominate grants** flip `self.GrantedGoAgain` /
  `self.GrantedDominate`; `EffectiveGoAgain` / `EffectiveDominate` honour the flag. Card
  docstrings call out the *condition*, not the flag.
- **`card.VariableCost` markers** (Amplify the Arknight, Rune Flash, …): `Cost(s)` reads
  TurnState; the marker exposes `MinCost` / `MaxCost` for the solver's pre-screen. Don't
  re-document the dispatch — note the printed cost formula.
- **Attack Reactions**: cards implement `sim.AttackReaction.ARTargetAllowed(c) bool`
  matching the printed target wording, and call `sim.GrantAttackReactionBuff(s, predicate,
  n)` from `Play` to add `+n{p}` to the first matching `CardsRemaining` entry. The partition
  validator rejects attack-role assignments where no chain card satisfies the predicate, so
  an AR with no legal target is unplayable rather than silently wasted. ARs cost 0 AP; the
  chain runner's free-step gate handles that automatically. Card docstrings call out the
  printed predicate (esp. when the wording distinguishes "attack" from "attack action card")
  and any modelling fudge — not the wiring.

## Logging idioms

Card.Play uses two orthogonal `TurnState` primitives: `AddValue(n)` mutates `s.Value`,
and `Log` / `LogRider` / `LogPreTrigger` / `LogPostTrigger` (plus their `f` formatted
variants) append a `LogEntry`. They never collude — each does one thing, and the
internal `skipLog` gate lives inside the Log helpers so cards never check it.

**A line that starts with `s.Log(...` (or any `Log*` helper) must have no side effects.**
Put the value change on its own preceding line:

```go
// Good — Log line is pure:
s.AddValue(s.CreateRunechants(2))
s.LogRider(self, 2, "Created 2 runechants")

// Bad — side effect hidden inside the Log call:
s.LogRider(self, s.AddValue(s.CreateRunechants(2)), "Created 2 runechants")
```

Reasons: a reader scanning for "what does this card do at runtime" can skip every
`Log*` line knowing it's just printout, and a future profile-driven optimisation that
short-circuits log construction (e.g. a `WantsLog()` gate cards consult locally) won't
silently drop the side effect.

For attack and defense chain steps the standard idiom is two lines: capture the credited
amount via `self.DealEffectiveAttack(s)` / `DealEffectiveDefense(s)` (which encapsulates
the AddValue side effect) and pass it to `s.Log(self, n)`:

```go
// Attack action:
n := self.DealEffectiveAttack(s)
s.Log(self, n)

// Defense reaction:
n := self.DealEffectiveDefense(s)
s.Log(self, n)

// Non-attack chain step:
s.Log(self, 0)
```

## Cross-file references

If a comment's rationale would otherwise cite "matches the pattern in foo.go, bar.go,
baz.go", factor the shared rule into this file and cite only the local behaviour at the call
site.

## Test layout

Two homes for tests:

- **Unit tests** live next to the code they cover (`internal/sim/foo_test.go` for
  `internal/sim/foo.go`, `internal/cards/foo_test.go` for `internal/cards/foo.go`, etc.).
  They may use the in-package `package sim` or the black-box `package sim_test` form. They
  may exercise unexported helpers via test exports — but only when no public entry point
  reaches the same behaviour. Each card under `internal/cards/` covers its own rider via a
  unit test that calls the card's `Play` directly.
- **End-to-end tests** live in the top-level `e2etest/` package. They exercise the
  simulator through public entry points only: `(*Deck).EvalOneTurnForTesting` for chain
  evaluation, `(*Deck).EvaluateWith` for full multi-turn runs. They use real heroes from
  `internal/heroes` (e.g. `heroes.Viserai{}`) rather than package-private stubs. Anything
  that would otherwise need an `exports_test.go` re-export goes here instead.

`sim.Best` and `sim.BestWithTriggers` carry a "Test convention" doc paragraph pointing at
`EvalOneTurnForTesting`. They aren't `// Deprecated:` because the simulator itself calls
them internally; the convention is for new test code only. New e2e tests should not call
`sim.Best` directly — they should drive the deck through `EvalOneTurnForTesting` so the
test mirrors production's per-turn loop.

### Test docstrings

A test's doc comment is a single brief sentence stating the behavior under test, e.g.
`// Tests that a single pitch paying for multiple Aether Slashes activates the bonus on
each.` Inputs, expected values, and the chain shape are visible in the test body and
don't belong in the comment. The same rule applies to unit tests and e2e tests.
