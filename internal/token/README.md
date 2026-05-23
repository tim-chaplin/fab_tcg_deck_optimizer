# internal/token

## Purpose

`token` owns the factories for Flesh and Blood's five built-in tokens: the item tokens
**Gold / Silver / Copper** and the aura-flavoured tokens **Runechant / Ponder**. Each factory
returns the concrete value the engine stores — `*item.Item` for item tokens, `*aura.Aura` for
aura tokens — pre-wired with the canonical name, the reserved token identifier, and (for the
aura tokens) the trigger type plus fire closure.

## Key functions

- `NewGold(n)` / `NewSilver(n)` / `NewCopper(n)` — token items at count `n`. Each wraps the
  matching card's activated ability so the chain runner can enqueue it as a playable.
- `NewRunechant(n)` — a `CardOrAbility`-triggered aura filtered to attacks. When an attack is
  played it flips `ArcaneDamageDealt` if its count clears the damage-likely-to-hit window,
  then destroys. Resolving before the attack's own effect lets it turn on that attack's
  "dealt arcane damage this turn" rider. Damage itself is credited at creation time inside
  `GameEngine.CreateRunechants`; the handler is pure state cleanup.
- `NewPonder(n)` — an `EndOfTurn`-triggered aura. For each token in play it pops the top of
  the deck into hand (letting the post-hoc arsenal-promotion step fill an empty arsenal),
  then destroys. Pops past deck-end are silently skipped.

## Key type

- **`GameEngine`** (`interfaces.go`) — the narrow slice of engine surface the Runechant fire
  closure needs beyond `card.GameEngine`. Runechant calls `SetArcaneDamageDealt`, which is not
  card-facing, so the closure type-asserts the engine to this interface.
  `*gameengine.GameEngine` satisfies it structurally.

## How it is used

`gameengine`'s token economy (`CreateGold`, `CreateRunechants`, …) calls these factories to
mint a fresh token entry, then `bumpOrCreateAura` / `bumpOrCreateItem` consolidates it: at
most one entry per token kind, with `Count` accumulating. Token entries are identified by
their canonical display name and a reserved `CardID` from `internal/ids`, so the eval-cache
key distinguishes each kind without a separate discriminator.

## How to extend it

A new built-in token adds a factory here (item or aura), a reserved ID in `internal/ids`, and
an engine `Create` / `Count` pair. Reuse the `bumpOrCreate` consolidation so the kind stays a
single counted entry.

## Important files

- `token.go` — the five factories.
- `interfaces.go` — the narrow `GameEngine` interface for the Runechant closure.

## Gotchas and invariants

- Aura-token fire closures here are constructed once per factory call and kept inline;
  per-variant payload rides through `Count`.
- Item tokens leave their trigger core zero-valued, so their trigger type never matches a
  firing event — they act only through their activated ability.
- Token destroy never routes the source to the graveyard (a token has no card).
- `CreateRunechants` credits damage at creation; the Runechant aura handler must not
  re-credit it.
</content>
