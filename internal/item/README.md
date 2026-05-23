# internal/item

## Purpose

`item` owns the concrete `Item` type — an in-play permanent the engine stores in its arena.
An item comes in two flavours sharing one struct: a **token item** carries an activated
ability the chain runner enqueues as a playable, and a **card-sourced item** carries a
trigger that fires on a scheduled event (mirroring an aura). Both embed the shared
`trigger.Trigger` core; token items leave it zero-valued so trigger type 0 never matches a
firing event.

## Key type

- **`Item`** — the concrete entry. Embeds `trigger.Trigger[card.Item]` for source/token
  identity and (for card-sourced items) the firing event and handler; adds the stack `count`
  and the optional activated `ability`.

## Key functions

- `NewFromToken(name, tokenID, ability, count)` — a token-sourced item. Production callers
  reach for `internal/token`'s `NewGold` / `NewSilver` / `NewCopper`.
- `NewFromCard(source, tt, fire, oncePerTurn, typeFilter)` — a card-sourced triggered item
  whose handler fires when an event of type `tt` resolves; `count` starts at 1; `typeFilter`
  narrows the firing site (pass `nil` for no filter).
- `(*Item).Fire` — invokes the stored handler with the item as typed receiver; sets
  `activeEngine` for the handler's duration so a handler-side `Destroy` routes back through
  `engine.DestroyItem`.
- `(*Item).Ability` — the activated-ability card boxed as `any` (callers assert to their
  richer card type), or `nil` for a card-sourced item with no ability.

## Item activated abilities

Items mirror the weapon split: a `sim.Item` permanent (`Self` / `Count` / `Ability`) lives in
`GameState.Items`, and a `card.Card` activated ability (e.g. `GoldTokenAbility`) is what the
chain runner enqueues each turn. The ability carries the parent's `card.TypeSet` (so
`TypeItem` keeps `PersistsInPlay` true and the chain step never hits the graveyard) plus any
subtype the printed text uses; an activated ability that does not attack omits `TypeAttack`.

Token items consolidate by token kind: at most one `Item` per kind per `GameState`. The
per-token `Create` helper (`g.CreateGold(n)`, …) bumps an existing entry's `Count` or appends
a new one. `g.ConsumeItem(t, n)` decrements `Count` and removes the entry at zero — the
standard "ability paid one charge" path. Token items don't head to the graveyard on destroy.

The chain runner builds the item-ability playable list by replicating each item's `Ability`
up to `perItemAbilityCap` times so the weapon/ability mask can pick "play it 0..N times this
turn"; the cap (in `internal/sim/sequence.go`) bounds the `2^k` mask explosion.

An item can instead carry a trigger. Because `Item` embeds the shared `trigger.Trigger` core,
a card-sourced item built by `NewFromCard` fires a handler on a scheduled event exactly like
an aura — the engine's `FireTriggers` walks items alongside auras, and a handler ends the
item's life via `DestroyItem`. Token items leave the trigger zero-valued, so their trigger
type never matches.

The `triggertype.Pitch` event fires as each card is pitched, with the pitched card as the
triggering card (read it via `GameEngine.TriggeringCard()`). A `Pitch` handler raises what
that pitch yields by calling `GameEngine.AddResourcePoints(n)`; the grant folds into the
pitched card's resource contribution. A card whose printed text reads "Whenever you pitch a
card, …" registers a `Pitch`-triggered item from its `Play` via
`GameEngine.CreateItem(self, triggertype.Pitch, handler, false, nil)` — Talisman of
Recompense is the model.

## Important file

- `item.go` — the entire package: the `Item` struct, the two constructors, `Fire`,
  `Destroy`, `Ability`, `Copy`.

## Gotchas and invariants

- `activeEngine` is live only for the duration of a `Fire`; calling `Destroy` outside a fire
  window panics deterministically.
- A token item's trigger core is zero-valued, so it never fires through `FireTriggers` — it
  acts only through its activated ability.
- `Item` exposes `Copy` (allocating) but no in-place `CopyInto` reset surface, so the
  engine's per-entry item copy always allocates; only the slice backing is pooled.
</content>
