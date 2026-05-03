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
  modelled" sentence in prose. Example:
  - **Demolition Crew** (Generic Action - Attack with Dominate + an additional reveal cost) —
    no "Modelling: Dominate is advertised via the `card.Dominator` marker..." block. The
    `Dominator` interface implementation makes that link by itself; the additional reveal cost
    is documented by the `// not implemented:` comment above its `NotImplemented` method.

## Aura lifecycle

Defined in `internal/sim/aura.go`. Standard shape for cards that "create an aura that fires
later":

- `Play` calls `s.AddAura(sim.Aura{...})` with `Self`, `TriggerType`, `Count`, and `Handler`.
  `AddAura` sets `s.AuraCreated = true` for same-turn aura-readers.
- The sim walks `s.Auras` on each matching `TriggerType` condition and invokes every matching
  `Handler`. `OncePerTurn` caps an Aura at a single fire per turn.
- **Lifecycle is the handler's job, not the sim's.** A handler that's done calls
  `s.DestroyAura(t, addToGraveyard)` to splice itself out of `s.Auras` and (when
  `addToGraveyard`) land `Self` in the graveyard. Counter-based auras decrement `t.Count` and
  call `DestroyAura` when it hits zero. The sim never mutates `Count` or graveyards on its own.
- Aura handlers parallel `Card.Play` — `func(s *sim.TurnState, t *sim.Aura)` with no return.
  Credit damage / life gain via `s.AddValue(n)`; emit log lines via `s.LogPostTrigger`.

**Handlers must be top-level functions, not inline closures inside `Play`.** Even closures that
capture nothing escape to the heap when assigned to `Aura.Handler` (`go build -gcflags='-m'`
confirms `func literal escapes to heap`), allocating one closure per Play. A top-level handler
is a static function pointer — zero alloc per registration. Pattern:

```go
func mySigilAuraHandler(s *sim.TurnState, t *sim.Aura) {
    s.AddValue(1)
    s.DestroyAura(t, true)
}

func (c MySigil) Play(s *sim.TurnState, self *sim.CardState) {
    s.AddAura(sim.Aura{Self: c, TriggerType: sim.TriggerStartOfTurn, Count: 1, Handler: mySigilAuraHandler})
    s.Log(self, 0)
}
```

Per-variant payloads (e.g. R/Y/B versions creating different rune counts) thread through
`Aura.Count` so the handler stays shared across variants.

Card docstrings should NOT restate this lifecycle. State only what's card-specific — the printed
clause and any rider the handler drops.

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
- **Plain-block bonuses** (Battlefront Bastion, Right Behind You, …): cards implement
  `sim.Blocker.Block(s, self)`. The chain runner calls Block on every plain blocker that
  opts in, with `s.Defenders` populated with the partition's full defender slice (DRs +
  plain blocks). Implementations scan Defenders and flip `self.BonusDefense` to credit
  conditional buffs ("+1{d} when defending alone", "+1{d} when defending together with
  another card", etc.); `defendersDamage` folds the printed Defense plus `BonusDefense`
  into the block, capped by remaining incoming damage. Cards without block-time logic
  don't implement Blocker; their plain-block contribution stays at the printed Defense.
- **Conditional go-again / dominate grants** flip `self.GrantedGoAgain` /
  `self.GrantedDominate`; `EffectiveGoAgain` / `EffectiveDominate` honour the flag. Card
  docstrings call out the *condition*, not the flag.
- **`card.VariableCost` markers** (Amplify the Arknight, Rune Flash, …): `Cost(s)` reads
  TurnState; the marker exposes `MinCost` / `MaxCost` for the solver's pre-screen. Don't
  re-document the dispatch — note the printed cost formula.
- **Attack Reactions**: cards implement `sim.AttackReaction.ARTargetAllowed(c, mode) bool`
  matching the printed target wording. The chain runner validates that the AR's chosen
  Mode accepts the active attack; failures abort the permutation as illegal. The AR's
  `Play` calls `sim.GrantAttackReactionBuff(s, self, n)` — the helper reads
  `s.AttackReactionTarget()` (set by the runner before invoking `Play`), adds `n` to the
  target's `BonusAttack`, credits `n` to `s.Value`, amends the buffed attack's chain-step
  display delta, and emits the rider log line. ARs cost 0 AP; the chain runner's free-step
  gate handles that automatically. Modal ARs combine this with `sim.ModalCard.Modes()` —
  the predicate dispatches per mode and `Play` applies the buff unconditionally because
  the runner already validated. Non-modal ARs ignore the `mode` parameter (always 0).
  Card docstrings call out the printed predicate (esp. when the wording distinguishes
  "attack" / "attack action card" / "weapon attack" — the corresponding TypeSet helpers
  are `IsAttack` / `IsAttackAction` / `IsWeaponAttack`) and any modelling fudge — not the
  wiring.
- **`OnHit` registrations**: attack cards with "if this hits, do X" riders append a
  `func(*sim.TurnState)` to `self.OnHit` inside `Play` instead of firing the rider inline.
  The chain runner finalizes each attack post-AR-buff: when `LikelyToHit(self)` evaluates
  true on the post-buff `EffectiveAttack`, every closure in `self.OnHit` runs. Cards that
  add an on-hit rider to a DIFFERENT card (Mauvrion Skies, Runic Reaping) append to the
  target's `OnHit`. Cards must NOT call `LikelyToHit` directly from `Play` — the chain
  runner owns the gate so AR buffs propagate.
- **`NextAttackActionHit` triggers** (Plunder Run, …): cards whose printed text reads "the
  next time an attack action card you control hits this turn, do X" register a
  `sim.NextAttackActionHitTrigger` via `s.RegisterNextAttackActionHit(...)` inside `Play`.
  The chain runner drains the queue inside `finalizeActiveAttack` on the first attack
  action that lands (`IsAttackAction` + `LikelyToHit`); all pending listeners fire together
  on that hit. Use this — not `OnHit` on a specific `CardState` — when the rider must wait
  across misses for the FIRST attack action that actually hits.
- **Self-granting on-hit go-again** (Overload, Razor Reflex mode 1): cards with a printed
  "if this hits, it gains go again" clause flip `self.GrantedGoAgain = true` inside `Play`
  when `sim.LikelyToHit(self)` returns true; `EffectiveGoAgain` honours the flag for the
  next chain step. Carry the `sim.ConditionalGoAgain` marker so the static lint test in
  `conditional_go_again_test.go` passes.
- **Modal "Choose 1" cards** (Captain's Call, …): cards implement `sim.ModalCard.Modes()
  int` returning the mode count and dispatch on `self.Mode` inside `Play`. The chain
  runner enumerates the cartesian product of mode indices across all modal cards in a
  permutation and picks the highest-Value tuple. Modes that are no-ops for the current
  state should resolve as zero-Value no-ops; the runner will pick a sibling mode that
  contributes more. Card docstrings call out each mode's effect — not the wiring.
- **Modal cost** (Bluster Buff / Chest Puff / Look Tuff cycle): ModalCards whose resource
  cost varies by mode implement `sim.ModalCost.ModalCost(mode int8) int`. The attacker
  meta cache folds the per-mode min/max into the partition pre-screen, and the chain
  runner reads the mode's cost via `costAt(state, self.Mode)` instead of the static /
  VariableCost path. Card docstrings call out the per-mode (cost, effect) pair.
- **Modal blockers** (Brothers in Arms, …): plain-block cards with mode-dependent
  block-time costs implement `sim.ModalCard.Modes()` + `sim.Blocker.Block(s, self)` +
  `sim.BlockCost(mode int8) int`. `defendersDamage` enumerates each modal blocker's
  modes within the partition's spare defense budget (`phase.defendBudget − drCost`) and
  picks the highest-BonusDefense mode that fits. Mode 0 is conventionally the printed
  default (cost 0, no extra effect). Card docstrings call out each mode's (cost, effect)
  pair — not the wiring.
- **`DefensiveInstant` markers** (Brush Off, Calming Breeze, Oasis Respite, Peace of
  Mind, …): `TypeInstant` cards whose printed effect prevents damage during the defense
  phase opt in via the `sim.DefensiveInstant` marker. The partition treats them as
  defenders, their `Cost()` is summed against the defense budget, and `Play` calls
  `self.DealEffectiveDefense(s)` so prevention is `min(Defense(), IncomingDamage)`.
  Damage prevention is collapsed against the sim's single `IncomingDamage` bucket: "the
  next N damage" and "the next K times … prevent 1 each" both reduce to a single
  `Defense() = N` bucket; "next damage of M or less" credits `min(M, IncomingDamage)`.
  Card docstrings note the printed prevention amount and any rider that's dropped — not
  the wiring or the bucketing rationale.

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
