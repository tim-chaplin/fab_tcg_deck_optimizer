# internal/gameengine

## Purpose

`gameengine` owns the per-turn game-state engine that the `sim` chain runner drives. It is
split into two cooperating types: **`GameState`** holds the raw per-turn data (every slice,
scalar, and flag), and **`GameEngine`** embeds a `*GameState` and adds the rules-engine API —
trigger dispatch, chain-step resolution, the token economy, deck manipulation, clash/opt.
Cards play against this engine through the narrow `card.GameEngine` interface, which is a
subset of `GameEngine`'s method set.

## The GameState / GameEngine split

`GameState` owns the data; `GameEngine` owns the rules. `*GameState` is embedded in
`GameEngine`, so every pure accessor and the `Copy` / `Reset` utilities promote
automatically. Internal machinery (sim's per-permutation scratch, carryover snapshots, the
winning pointer) passes around a bare `*GameState` when it only needs to read or copy data;
the chain runner wraps it in a `*GameEngine` via `GameState.Engine()` to drive `Card.Play`.
`GameEngine` methods that shadow a promoted `GameState` method add cacheable-flipping or
rules logic on top.

## Key types and interfaces

- **`GameState`** (`state.go`) — the raw state. Construct via `GameStateBuilder()`; copy via
  `Copy` (deep), `CopyFrom` (in-place, pool-reusing), or `CopyPersistentState` /
  `CopyPersistentStateFrom` (carryover-only). `ResetEphemeralState` returns it to the
  start-of-turn baseline.
- **`GameEngine`** (`engine.go`) — the rules wrapper. `New()` returns one over a default
  state.
- **`Aura` / `EphemeralTrigger` / `Item` / `Hero`** (`interfaces.go`) — the engine's local
  interfaces. Concrete impls live in leaf packages (`internal/aura`, `internal/trigger`,
  `internal/item`, `internal/hero`, `internal/token`) that do **not** import `gameengine`;
  `sim` wires builders in via `init` so the engine never imports the concrete types. `Hero`
  carries the trigger-walking surface (`TriggerType` / `Fire` / etc.) so `FireTriggers`
  dispatches the hero alongside auras / items.
- **`StateBuilder`** (`builder.go`) — fluent `GameState` construction. Starts from sensible
  defaults (cacheable, no-op logger, empty deck, a 20-health/4-intelligence fallback hero).

## How it is used / how to extend it

The `sim` chain runner builds one master `*GameState`, copies it per partition leaf, runs the
defense pass, then runs each chain permutation against a fresh per-permutation copy. Cards
reach the engine through hooks (`Play`, `Block`, `OnHitHandler.Fire`, trigger handlers,
`Hero.OnCardPlayed`), all of which receive the `card.GameEngine` / `card.Logger` typed
surfaces.

Adding a new matchup-level parameter or per-turn flag adds a field on `GameState` plus its
accessor pair, and — if cards mutate it from hidden state — a cacheable-flipping shadow on
`GameEngine`. A new lifecycle event adds a `triggertype.Type` bit and a `FireTriggers` call
site.

## Trigger dispatch

`FireTriggers(t, triggeringCard)` is the single dispatch point for every lifecycle event. It
fires the hero (via the standalone `fireHero` helper) and then walks auras, ephemeral
triggers, and items via `fireHooks`. `fireHooks` takes a length snapshot up front (so an
entry a handler creates lands past it and is not fired this pass — the self-exclusion
mechanism); honours the once-per-turn gate; applies each entry's type filter against the
triggering card's type set; and uses a cursor walk so a handler-side destroy doesn't skip
the next entry. Ephemeral triggers are spliced out unconditionally after firing
(one-shot).

The hero is one more triggered entity but singular — no slice splicing, no removeAfterFire
— so `fireHero` applies the OncePerTurn / Matches gates directly without the cursor walk.
A triggered hero (`internal/hero/heroes`) embeds `trigger.Trigger[card.Hero]` to provide
the trigger surface; heroes with no ability (`defaultHero`, test stubs) return
`TriggerType == 0` so the dispatch's bit-and check skips them.

## Chain-step resolution

`ResolveChainStep(l, pc)` runs `pc.Card.Play`, then `chainStepDelta` credits the standard
chain-step value to `value` and `AppendChainStep` appends the canonical
`<DisplayName>: <VERB> (+N)` log entry — appended *after* `Play` so self-buffs are reflected.
Attacks credit `EffectiveAttack`; defense reactions and `DefensiveInstant` cards credit
`EffectiveDefense` capped at the remaining unblocked damage; everything else logs `(+0)`.

## Important files

- `state.go` — `GameState`, the copy/reset family, all pure accessors.
- `engine.go` — `GameEngine`: cacheable-flipping zone accessors, `FireTriggers` / `fireHooks`,
  `DestroyAura` / `DestroyItem`, `ResolveChainStep`, `Opt` / `Clash`, arcane damage, the
  token economy.
- `interfaces.go` — the `Aura` / `EphemeralTrigger` / `Item` / `Hero` interfaces.
- `builder.go` — `StateBuilder` and `New()`.
- `aura.go` — aura-list helpers bridging the engine to `internal/aura`.
- `heuristics.go` — `LikelyToHit` / `LikelyDamageHits` and damage-value constants.
- `logger.go` — `NoopLogger` and the logger plumbing.
- `default_hero.go` — the no-ability fallback hero.

## Gotchas and invariants

- Card-facing zone accessors (`Hand`, `Graveyard`, `Deck`, `PeekTopN`, …) flip `cacheable`
  to false as a side effect of exposing hidden state. The engine's own internals use the
  non-flipping promoted `GameState` variants (`ge.GameState.X`).
- The hand is kept sorted by `Card.ID()` on every insert (`insertHandSorted`) so it stays a
  canonical multiset for the chain runner and the eval-cache key.
- `ResetEphemeralState` keeps cross-turn carryover (hero, deck, hand, arsenal, graveyard,
  banished, auras, items, opponentMarked, incoming damage) and wipes everything else;
  `incomingDamage` is the constant matchup figure and survives, but `damageBlocked` resets.
- `Copy` resets the logger to nil/`NoopLogger` — callers install a fresh per-clone logger
  only when recording, keeping find-best copies allocation-free.
- `Opt` panics if the hero handler's `(top, bottom)` output is not exactly the input
  multiset.
</content>
