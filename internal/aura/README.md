# internal/aura

## Purpose

`aura` owns the concrete `Aura` type — a persistent hook entry that the game engine stores in
its arena and fires at a scheduled trigger type (start of turn, when a card is played, end of
turn, …). Card-backed and token-backed auras share the same struct; `SourceCard`
distinguishes them. An aura is the modelling primitive for any card that "creates something
that fires later".

## Key types

- **`Aura`** — the concrete entry. Embeds `trigger.Trigger[card.Aura]` for the shared core
  (trigger type, handler, source identity, optional type filter, once-per-turn gate) and adds
  a multi-fire `count`.
- **`Handler`** — the typed handler signature: `func(card.GameEngine, card.Logger, card.Aura)`.

## Key functions

- `NewFromCard(source, tt, fire, count, oncePerTurn, typeFilter)` — a card-backed aura.
  `SourceCard` surfaces the originating card so the engine can route it to the graveyard on
  destroy.
- `NewFromToken(name, tokenID, tt, fire, count, typeFilter)` — a token aura with no
  originating card; `CardID` returns the token ID so cache keys distinguish each token kind.
- `(*Aura).Fire` — invokes the stored handler with the aura as the typed receiver. It sets
  `activeEngine` for the handler's duration so a handler-side `Destroy` can route back
  without allocating a per-fire wrapper.
- `(*Aura).Copy` / `CopyInto` — deep copy boxed as `any`; `CopyInto` rewrites a pooled prior
  slot in place for the per-permutation reset fast path.

## Aura lifecycle

For cards that create an aura which fires later:

- `Play` calls `g.AddAura(...)` (or a token-aura helper) with the source card, trigger type,
  count, and handler. Aura creation sets `g.AuraCreated = true` for same-turn aura-readers.
- The engine walks its aura list per matching trigger type and invokes every matching
  handler. `OncePerTurn` caps an aura at one fire per turn; the gate is re-armed at the turn
  boundary.
- **Lifecycle is the handler's job.** When done, the handler calls
  `ctx.Destroy(addToGraveyard)` — which routes through the firing engine's `DestroyAura`, to
  splice the aura out of the arena and, when `addToGraveyard`, land `SourceCard` in the
  graveyard. Counter-based auras decrement `Count` and destroy at zero. The engine never
  mutates `Count` or graveyards on its own.
- Handlers parallel `Card.Play`: `func(g card.GameEngine, l card.Logger, ctx card.Aura)`, no
  return. Credit damage / life via `g.AddValue(n)`; emit log lines via the post-trigger
  logger methods.

The engine-side walk that fires aura handlers (`FireTriggers` / `fireHooks`) lives in
`internal/gameengine` — see `internal/sim` and `internal/gameengine` for how the attack-turn runner
and start-of-turn pass invoke it. `fireHooks` takes a length snapshot before firing so an
aura created by a handler is not consumed in the same pass, and uses a cursor walk so a
handler-side destroy doesn't skip the next entry.

**Card-backed aura handlers must be top-level functions, not inline closures.** A closure
passed to `NewFromCard` escapes to the heap — one allocation per `Play` — whereas a top-level
handler is a static function pointer. Built-in token auras are exempt: their factories in
`internal/token` construct the closure once at factory time, not per `Play`. Per-variant
payloads (R/Y/B rune counts, etc.) thread through `Aura.Count` so a single shared handler
covers every variant:

```go
func mySigilAuraHandler(g card.GameEngine, l card.Logger, ctx card.Aura) {
    g.AddValue(1)
    ctx.Destroy(true)
}
```

## Important file

- `aura.go` — the entire package: the `Aura` struct, the two constructors, `Fire`, the copy
  family, and `Destroy`.

## Gotchas and invariants

- `activeEngine` is set only for the duration of a `Fire`. Calling `Destroy` outside a fire
  window panics deterministically rather than silently no-oping.
- `Copy` / `CopyInto` clear `activeEngine` on the copy so a stale engine pointer can't leak
  across permutations.
- Per-permutation `Count` / `FiredThisTurn` mutations stay isolated because the engine
  deep-copies each aura when it copies a `GameState`.
</content>
