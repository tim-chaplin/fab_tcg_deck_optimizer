# Developer Standards

Shared conventions; rules cited across multiple files factor here so per-file comments stay card-specific.

## Card file layout

Cards live in `internal/cards/`. A card file typically:

1. Declares a shared `TypeSet` var for the card's type line.
2. Declares one struct per printed pitch variant (`FooRed`, `FooYellow`, `FooBlue`).
3. Implements `card.Card` plus optional markers (`VariableCost`, …) on each variant.
4. Shares a `fooPlay(...)` helper when variants differ only by a numeric parameter.

Card data (name, cost, pitch, attack, defense, type line, printed text) is transcribed from `github.com/the-fab-cube/flesh-and-blood-cards`; no per-file `Source:` line.

## Comment rules

- Wrap at 100 chars.
- Describe CURRENT behavior + motivation. No history references ("replaces X", "was Y before", "now does Z"); no "previously / formerly / legacy / deprecated" framing.
- Delete completed TODOs instead of rewording them.
- Card docstrings cover card-SPECIFIC quirks — the printed rules text, any modelling fudge, and anything surprising about how this card interacts with the solver. Generic framework plumbing belongs in framework docs, not repeated per card.
- Don't restate behavior already documented by an external function, type, or marker the card uses. A card carrying `card.Dominator` doesn't re-explain Dominate; a card with `NotImplemented` + a `// not implemented: <quirk>` line doesn't repeat "rider isn't modelled" in prose.

## Aura lifecycle

Defined in `internal/sim/aura.go`. For cards that "create an aura that fires later":

- `Play` calls `s.AddAura(sim.Aura{...})` with `Self`, `TriggerType`, `Count`, `Handler`. `AddAura` sets `s.AuraCreated = true` for same-turn aura-readers.
- Sim walks `s.Auras` per matching `TriggerType` and invokes every matching `Handler`. `OncePerTurn` caps an Aura at one fire per turn.
- **Lifecycle is the handler's job.** When done, the handler calls `s.DestroyAura(a, addToGraveyard)` to splice itself out of `s.Auras` and (when `addToGraveyard`) land `Self` in the graveyard. Counter-based auras decrement `a.Count` and call `DestroyAura` at zero. The sim never mutates `Count` or graveyards on its own.
- Handlers parallel `Card.Play` — `func(s *sim.TurnState, t *sim.Trigger, a *sim.Aura)`, no return. Aura fires pass the firing aura as `a`; standalone trigger fires pass `nil`. Credit damage / life via `s.AddValue(n)`; emit log via `s.LogPostTrigger`.

**Handlers must be top-level functions, not inline closures.** Closures assigned to `Aura.Handler` escape to the heap (`go build -gcflags='-m'` shows `func literal escapes to heap`), allocating one per Play; top-level handlers are static function pointers. Per-variant payloads (R/Y/B rune counts) thread through `Aura.Count` so the handler stays shared:

```go
func mySigilAuraHandler(s *sim.TurnState, _ *sim.Trigger, a *sim.Aura) {
    s.AddValue(1)
    s.DestroyAura(a, true)
}

func (c MySigil) Play(s *sim.TurnState, self *sim.CardState) {
    s.AddAura(sim.Aura{Self: c, TriggerType: sim.TriggerStartOfTurn, Count: 1, Handler: mySigilAuraHandler})
    s.Log(self, 0)
}
```

Card docstrings state only what's card-specific — the printed clause and any rider the handler drops.

## NotImplemented vs Unplayable markers

Both exclude a card from random / mutation pools:

- `sim.NotImplemented` — placeholder; card *would* be worth modelling. Pair with a `// not implemented: <one-line description>` comment above the `NotImplemented()` method. Files in `internal/cards/notimplemented/`.
- `sim.Unplayable` — verdict; card too weak to want even if modelled. **No per-card rationale in the docstring.** Files in `internal/cards/unplayable/`.

`ls internal/cards/notimplemented/` is the live todo list. `TestLayout_MarkersStayInSubpackages` enforces the layout.

## Standard rider wiring

The plumbing below is uniform; card docstrings call out the printed rider and any modelling fudge, never the wiring.

- **Played-from-arsenal go-again** (Fervent Forerunner, Frontline Scout, Performance Bonus, Promise of Plenty, Scour the Battlescape, …): `self.GrantGoAgainIfFromArsenal()` at top of Play. Flips `GrantedGoAgain` only when this copy came from the arsenal slot.
- **+N{d} on arsenal-played defense reactions** (Springboard Somersault, Unmovable, …): implement `card.ArsenalDefenseBonus` returning `N`; `CardState.EffectiveDefense` folds it in for the arsenal-in copy.
- **Plain-block bonuses** (Battlefront Bastion, Right Behind You, …): implement `sim.Blocker.Block(s, self)`. Chain runner calls `Block` on every plain blocker that opts in, with `s.Defenders` set to the partition's full defender slice (DRs + plain blocks). Implementations scan `Defenders` and flip `self.BonusDefense` for conditional buffs ("+1{d} when defending alone", "+1{d} together with another card", …); `defendersDamage` folds printed Defense + `BonusDefense`, capped by remaining incoming.
- **Conditional go-again / dominate grants**: flip `self.GrantedGoAgain` / `self.GrantedDominate`; `EffectiveGoAgain` / `EffectiveDominate` honour them. Card docstrings call out the *condition*, not the flag.
- **`card.VariableCost`** (Amplify the Arknight, Rune Flash, …): `Cost(s)` reads TurnState; the marker exposes `MinCost` / `MaxCost` for the solver pre-screen. Note the printed cost formula.
- **Attack Reactions**: implement `sim.AttackReaction.ARTargetAllowed(c, mode) bool` matching the printed target wording. Chain runner validates the chosen `Mode`; failures abort the permutation. The AR's `Play` calls `sim.GrantAttackReactionBuff(s, self, n)` — reads `s.AttackReactionTarget()` (set by the runner before `Play`), adds `n` to target's `BonusAttack`, credits `n` to `s.Value`, amends the buffed attack's chain-step delta, emits the rider line. ARs cost 0 AP (handled by the free-step gate). Modal ARs combine with `sim.ModalCard.Modes()` — predicate dispatches per mode; `Play` applies the buff unconditionally because the runner already validated. Non-modal ARs ignore `mode` (always 0). Card docstrings call out the printed predicate (esp. when wording distinguishes "attack" / "attack action card" / "weapon attack" — TypeSet helpers are `IsAttack` / `IsAttackAction` / `IsWeaponAttack`) and any modelling fudge.
- **`OnHit` registrations**: attack cards with "if this hits, do X" append a `func(*sim.TurnState)` to `self.OnHit` inside `Play`. Chain runner finalizes each attack post-AR-buff: when `LikelyToHit(self)` is true on post-buff `EffectiveAttack`, every closure in `self.OnHit` runs. Cards adding an on-hit rider to a DIFFERENT card (Mauvrion Skies, Runic Reaping) append to the target's `OnHit`. **Do NOT call `LikelyToHit` directly from `Play`** — chain runner owns the gate so AR buffs propagate. Use a top-level handler (not an inline closure) so registration stays allocation-free — closures assigned to `OnHit.Fire` escape to the heap.
- **`NextHit` triggers** (Plunder Run, High Striker, …): cards whose printed text reads "the next time an X you control hits this turn, do Y" register via `s.RegisterNextHit(...)`. Each trigger carries a `TypeFilter func(card.TypeSet) bool` narrowing the qualifying hits — `card.TypeSet.IsAttackAction` for "attack action card" wording (Plunder Run), `card.TypeSet.IsAttack` for the broader "attack" wording that includes weapon swings (High Striker). Chain runner drains matching triggers inside `finalizeActiveAttack` on each `LikelyToHit` attack; non-matching triggers stay queued. Use this — not `OnHit` on a specific `CardState` — when the rider must wait across misses for the FIRST qualifying hit.
- **Self-granting on-hit go-again** (Overload, Razor Reflex mode 1): "if this hits, it gains go again" → flip `self.GrantedGoAgain = true` inside `Play` when `sim.LikelyToHit(self)` returns true.
- **Modal "Choose 1" cards** (Captain's Call, …): implement `sim.ModalCard.Modes() int` and dispatch on `self.Mode` inside `Play`. Chain runner enumerates the cartesian product of modes across modal cards and picks the highest-Value tuple. Modes that are no-ops should resolve as zero-Value no-ops. Card docstrings call out each mode's effect.
- **Modal cost** (Bluster Buff / Chest Puff / Look Tuff cycle): ModalCards whose resource cost varies by mode implement `sim.ModalCost.ModalCost(mode int8) int`. Attacker meta cache folds per-mode min/max into the partition pre-screen; chain runner reads the mode's cost via `costAt(state, self.Mode)`. Card docstrings call out each (cost, effect) pair.
- **Modal blockers** (Brothers in Arms, …): plain blockers with mode-dependent block-time cost implement `sim.ModalCard.Modes()` + `sim.Blocker.Block(s, self)` + `sim.BlockCost(mode int8) int`. `defendersDamage` enumerates modes within spare defense budget (`phase.defendBudget − drCost`) and picks the highest-`BonusDefense` mode that fits. Mode 0 is conventionally the printed default (cost 0, no extra effect). Card docstrings call out each (cost, effect) pair.
- **`DefensiveInstant`** (Brush Off, Calming Breeze, Oasis Respite, Peace of Mind, …): `TypeInstant` cards whose printed effect prevents damage opt in via `sim.DefensiveInstant`. Partition treats them as defenders; `Cost()` is summed against the defense budget; `Play` calls `self.DealEffectiveDefense(s)` so prevention is `min(Defense(), IncomingDamage)`. Damage prevention collapses against the single `IncomingDamage` bucket: "the next N damage" and "the next K times … prevent 1 each" both reduce to `Defense() = N`; "next damage of M or less" credits `min(M, IncomingDamage)`. Card docstrings note the printed prevention amount and any rider that's dropped.
- **Optional additional costs** (Looking for a Scrap, Nimble Strike, Regurgitating Slog, …): "you may pay X" gates where the cost spends graveyard residency (or other state the sim doesn't otherwise value) for a strictly-upside buff are paid unconditionally; the sim doesn't enumerate the skip branch. Card docstrings note the printed gate; no per-card rationale for "always tries".
## Logging idioms

`Card.Play` uses two orthogonal `TurnState` primitives: `AddValue(n)` mutates `s.Value`; `Log` / `LogRider` / `LogPreTrigger` / `LogPostTrigger` (plus `f` variants) append a `LogEntry`. The `skipLog` gate lives inside the Log helpers — cards never check it.

**A line starting with `s.Log(...` (or any `Log*` helper) must have no side effects.** Put the value change on its own preceding line:

```go
// Good — Log line is pure:
s.AddValue(s.CreateRunechants(2))
s.LogRider(self, 2, "Created 2 runechants")

// Bad — side effect hidden inside Log:
s.LogRider(self, s.AddValue(s.CreateRunechants(2)), "Created 2 runechants")
```

A reader scanning for "what does this card do" can skip every `Log*` line; a future profile-driven optimisation that short-circuits log construction won't silently drop the side effect.

For attack and defense chain steps the standard idiom is two lines — capture the credited amount via `self.DealEffectiveAttack(s)` / `DealEffectiveDefense(s)` (which encapsulates `AddValue`) and pass it to `s.Log(self, n)`:

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

## Weapon activated abilities

Weapons live in `internal/weapons/` as two paired Go types:

- A `sim.Weapon` permanent (`ID`, `Name`, `Types`, `Hands`, `Ability`) — the equipped permanent that sits in the arena and never enters the chain.
- A `sim.Card` activated ability (`<Weapon>Ability`) — what the chain runner enqueues when the player swings. Carries the printed Cost / Pitch / Attack / Defense / GoAgain / Play and the parent's `card.TypeSet` plus `card.TypeAttack` (so `IsAttack`, `IsWeaponAttack`, and — for Runeblade weapons — `IsRunebladeAttack` all fire).

`Ability()` returns a package-level cached `sim.Card` (each weapon file declares `var <weapon>Ability sim.Card = <Weapon>Ability{}` and returns it). The chain runner only calls `Ability()` once per weapon at attackBufs construction today, but the cache keeps any future hot-loop caller alloc-free since Go's interface boxing of zero-size structs allocates per call when the result escapes. The ability `Name()` matches the weapon's so chain logs read `<Weapon>: WEAPON ATTACK`.

IDs: weapon permanents take `WeaponID`; abilities take `CardID` and anchor in `internal/registry/ids/weapon_ids.go` past the last weapon-permanent ID. Token activated abilities (e.g. `GoldTokenAbilityID`) anchor past the last weapon-ability ID in the same file. Test fakes anchor past the last token-ability ID.

Card-attack predicates (`internal/card/types.go`) gate purely on `TypeAttack` (and `TypeWeapon` / `TypeRuneblade` where the rider needs the subtype). A bare weapon permanent — types only, no `TypeAttack` — never matches.

## Item activated abilities

Items mirror the Weapon split: a `sim.Item` permanent (Self / Count / Ability) lives in `TurnState.Items`, and a `sim.Card` activated ability (e.g. `GoldTokenAbility`) is what the chain runner enqueues each turn. The Ability carries the parent's `card.TypeSet` (so `TypeItem` keeps `PersistsInPlay` true and the chain step never hits the graveyard) plus any subtype the printed text uses; activated abilities that don't attack omit `TypeAttack`.

Token items consolidate by `TokenType`: at most one `Item` per `TokenType` per `TurnState`. The per-token `Create` helper (`s.CreateGold(n)`) bumps an existing entry or appends a new one. `s.ConsumeItem(t, n)` decrements `Count` and removes the entry at zero — the standard "ability paid one charge" path. Token items don't head to the graveyard on destroy.

The chain runner builds `ctx.itemAbilities` by replicating each Item's `Ability` up to `perItemAbilityCap` times so the wmask can pick "play it 0..N times this turn"; the cap (`internal/sim/sequence.go`) bounds the 2^k mask explosion.

## Cross-file references

If a comment's rationale would otherwise cite "matches the pattern in foo.go, bar.go, baz.go", factor the shared rule into this file and cite only the local behaviour at the call site.

## Test layout

Two homes:

- **Unit tests** next to the code (`internal/sim/foo_test.go` for `internal/sim/foo.go`, `internal/cards/foo_test.go` for `internal/cards/foo.go`). May use `package sim` or black-box `package sim_test`. May exercise unexported helpers via test exports — but only when no public entry point reaches the same behaviour. Each card under `internal/cards/` covers its own rider via a unit test calling `Play` directly.
- **Turn-level tests** in top-level `turntests/`. Public entry points only: `(*Deck).EvalOneTurnForTesting` for chain evaluation, `(*Evaluator).Evaluate` for full multi-turn runs. Use real heroes from `internal/heroes` (e.g. `heroes.Viserai{}`) rather than package-private stubs. Anything that would otherwise need an `exports_test.go` re-export goes here.

`sim.Best` and `sim.BestWithTriggers` carry a "Test convention" doc paragraph pointing at `EvalOneTurnForTesting`. Not `// Deprecated:` because the simulator itself calls them; the convention is for new test code only.

### Test docstrings

A test's doc comment is a single brief sentence stating the behavior under test, e.g. `// Tests that a single pitch paying for multiple Aether Slashes activates the bonus on each.` Inputs, expected values, and chain shape are visible in the test body. Same rule for unit and turn-level tests.
