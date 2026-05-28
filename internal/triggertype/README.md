# internal/triggertype

## Purpose

`triggertype` categorises *when* an aura or one-shot trigger fires. It is a deliberately
tiny leaf package — a single enum — so any consumer can name a lifecycle event without
pulling in the whole game engine.

## Key type

- **`Type`** — a bitmask of lifecycle events. Each constant is one bit. An aura or trigger
  subscribes to one event or an OR of several, and the engine's `FireTriggers` fires it when
  the dispatched event's bit is set.

## The events

- `StartOfTurn` — start of the owning player's action phase, before the best-line search.
- `CardOrAbility` — fires once as a card or weapon attack is played during the attack turn, before
  that card's own effect resolves. Subscribers narrow with a `TypeFilter` (e.g. attacks only
  for Runechant tokens, attack-actions for Malefic Incantation).
- `EndOfTurn` — after the attack turn finishes resolving, before the carry snapshot.
- `Hit` — when an attack hits (post-AR-buff `EffectiveAttack` survives blocks).
- `DamageTaken` — end of the defense phase, when incoming damage got through unblocked.
- `Pitch` — as each card is pitched to fund a play in the attack phase. The triggering card
  is the pitched card; a handler may boost the resources it yields via
  `GameEngine.AddResourcePoints`.
- `CrowdCheer` / `CrowdBoo` — when the crowd cheers / boos your hero, raised by
  `GameState.CrowdCheer` / `GameState.CrowdBoo`.
- `AttackBuffedByReaction` — when an attack action card's `BonusAttack` increments by
  exactly 1 during the reaction-step buff resolution. The buffed attack `CardState` is the
  firing context. The "exactly +1" gate lives at the fire site (in
  `CardState.GrantAttackReactionBuff`), not in subscribers.

## How it is used / how to extend it

`internal/trigger`, `internal/aura`, `internal/item`, and `internal/gameengine` all import
this enum to type their trigger entries and dispatch calls. A new lifecycle event adds one
`iota` constant here and a `FireTriggers` call site in the engine at the point the event
occurs.

## Important file

- `triggertype.go` — the `Type` enum and its lifecycle-event constants.

## Gotchas and invariants

- The constants are bit flags (`1 << iota`), so an entry can subscribe to several events with
  an OR; the engine's dispatch is a bitwise-AND test. Keep new constants single-bit.
- The package imports nothing — keep it dependency-free so the whole codebase can name a
  trigger type cheaply.
