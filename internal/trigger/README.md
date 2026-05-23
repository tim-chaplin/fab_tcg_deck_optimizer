# internal/trigger

## Purpose

`trigger` owns the shared trigger machinery: the embeddable `Trigger[T]` core that every
triggered entry carries, and the one-shot `EphemeralTrigger` kind. A trigger pairs a firing
event with a typed handler, the source identity, and an optional type filter. Auras
(`internal/aura`), items (`internal/item`), and triggered heroes (`internal/hero/heroes`)
embed `Trigger[T]` for their shared behaviour; `EphemeralTrigger` is the standalone
one-shot listener the engine fires once and drops.

## Key types

- **`Trigger[T]`** — the generic embeddable core. `T` is the concrete surface the handler
  receives (`card.EphemeralTrigger`, `card.Aura`, or `card.Item`), so handlers stay typed
  with no assertion. Holds the trigger type, the fire func, the source card (or token name +
  ID), the type filter, and the once-per-turn gate. The embedding type supplies a `Fire`
  method that calls `Invoke` with itself as the typed receiver.
- **`EphemeralTrigger`** — a one-shot card-sourced trigger. Carries nothing beyond the
  embedded core. The engine fires it on the next matching event and removes it from the
  queue.
- **`TypeFilter`** — `func(card.TypeSet) bool`, narrowing the firing site to a card-type
  predicate. `nil` means any matching event qualifies.

## Key functions

- `FromCard[T](source, tt, fire, oncePerTurn, typeFilter)` — a card-sourced core. `CardName`
  / `CardID` resolve from the source card.
- `FromToken[T](name, tokenID, tt, fire, oncePerTurn, typeFilter)` — a token-sourced core
  with no originating card; `CardName` / `CardID` return the supplied token identity.
- `FromHero[T](tt, fire, oncePerTurn, typeFilter)` — a hero-sourced core with no card
  or token identity. The embedding `Hero` type owns its own `ID` / `Name`, so the source /
  token slots stay zero and `CardName` / `CardID` return zero values on this core.
- `NewEphemeralTrigger(source, tt, fire, typeFilter)` — a one-shot card-sourced trigger.

## How it is used / how to extend it

The engine drives ephemeral triggers via `GameEngine.AddTrigger`, which wraps a card's
handler in an `EphemeralTrigger` and queues it. `FireTriggers` walks the queue per event
and splices out every fired entry.

A new triggered-entry kind (beyond aura / item / ephemeral) embeds `Trigger[T]` for the
chosen handler surface and supplies a `Fire` method calling `Invoke`. A new firing event adds
a `triggertype.Type` bit and a dispatch site in the engine.

## Important files

- `trigger.go` — `Trigger[T]`, `TypeFilter`, the `FromCard` / `FromToken` constructors,
  the accessors (`TriggerType`, `Matches`, `CardName`, `CardID`, `SourceCard`, `Invoke`).
- `ephemeral_trigger.go` — `EphemeralTrigger` and `NewEphemeralTrigger`.

## Gotchas and invariants

- `EphemeralTrigger` is one-shot: the engine drops it after firing, so its `oncePerTurn` is
  always false and its `FiredThisTurn` gate is never consulted — the trio of gate methods
  exists only to satisfy the shared `triggerHook` interface the engine's `fireHooks` walks.
- `Matches` returns true when the type filter is `nil`; turn-boundary events with no
  triggering card skip the filter entirely.
- `Trigger[T]` is effectively immutable after construction (only the `firedThisTurn` gate
  changes), so the engine duplicates only the slice header — not the entries — when copying
  the trigger queue.
</content>
