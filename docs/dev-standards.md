# Developer Standards

Shared conventions for code and comments. Rules cited across multiple files factor in here so
the per-file comments can stay card-specific.

## Card file layout

Each card implementation lives in `internal/card/generic/` or `internal/card/runeblade/` (plus
`internal/card/fake/` for test doubles). A card file typically:

1. Declares a shared `TypeSet` var for the card's type line.
2. Declares one struct per printed pitch variant (e.g. `FooRed`, `FooYellow`, `FooBlue`).
3. Implements the `card.Card` interface plus any optional markers (`NoMemo`, `AddsFutureValue`,
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
  (how `AuraTrigger` is ticked, what `NoMemo` opts out of, how memoization works) belongs in
  framework docs in `internal/card/card.go` and `internal/hand/hand.go`, not repeated in every
  card.
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

## NoMemo

Defined in `internal/card/card.go`. Cards whose `Play` depends on state not captured by the memo
key (deck composition, graveyard contents, cross-turn trigger state) implement `NoMemo()` so
`hand.Best` skips the cache. Restate only the card-specific reason in the docstring if it's not
obvious.

## Memoization

`hand.Best` memoizes by `(heroID, sorted weapon IDs, sorted card IDs, incomingDamage,
runechantCarryover, arsenal-in ID)`. Non-nil `priorAuraTriggers` or any `NoMemo` card in the
hand disables the cache for that call.

## Cross-file references

If a comment's rationale would otherwise cite "matches the pattern in foo.go, bar.go,
baz.go", factor the shared rule into this file and cite only the local behaviour at the call
site.
