# Mid-turn card draw: implementation plan

## Status

Active. Scratch pad for coordinating the phased rollout. Delete this doc
in the final commit before the PR goes out.

## Current model (baseline)

- A hand is up to N cards (typically 4). The solver enumerates every
  `{Pitch, Attack, Defend, Held, Arsenal}^N` partition, checks legality
  (cost coverage, phase feasibility), computes value, and takes the max.
- "Draw a card" riders are credited as a flat `card.DrawValue` (=3)
  regardless of turn state or top-of-deck identity. Cards that do this
  today: Snatch (on-hit), Drawn to the Dark Dimension (on-play), Sigil
  of the Arknight (conditional on top-of-deck type).

## Phased rollout

Each phase extends the set of dispositions the solver is allowed to
assign to a mid-turn-drawn card. The first phase removes the flat credit
and locks everything to HELD; subsequent phases unlock one disposition at
a time. This mirrors the five legal options for the drawn card in real
play: HELD, ARSENAL, PLAY, PITCH, DEFEND.

### Phase 0 — HELD (replaces the flat credit)

Baseline behaviour change:

- Introduce `state.DrawOne()`, a single shared helper on `TurnState`
  that advances `state.Deck` by one and appends the drawn card to a new
  `state.Drawn` slice. All cards with a "draw a card" rider (Snatch,
  Drawn to the Dark Dimension) call this once per draw event.
- Drop the `+ card.DrawValue` credit from those cards' Play methods.
- Model the drawn card as HELD: it displaces one card that would have
  been pulled at the end-of-turn refill step, so the net card count is
  unchanged and the credit is 0.

Why this is the right anchor: under the HELD-only model the simulator
neither over-credits (no imagined future tempo) nor under-credits
(we're not pretending the draw vanished — the card is tracked, just not
used productively yet).

**Sigil of the Arknight is out of scope for phase 0.** Its rider fires
at the start of next turn's action phase, *after* the end-of-turn
refill has drawn 4 cards, so the card is additive rather than
displacing a refill draw. Leaves its current flat-credit path alone
until a later phase models cross-turn additive draws explicitly.

Test coverage: end-to-end hand tests that exercise a draw-rider card
(e.g. Snatch) and pin the expected `TurnSummary.Value` at the new
(lower) number. Pre-phase-0 tests that asserted the `+3` credit move to
regression cases explicitly documenting the old behaviour's wrongness.

### Phase 1 — ARSENAL

Add ARSENAL as a disposition for drawn cards. Unlike HELD, arsenaling
does not displace the end-of-turn refill (arsenal is a separate zone),
so it's a positive-tempo option.

- The partition enumerator learns to consider `Drawn[k].Role = Arsenal`
  alongside the existing roles. Only one card can occupy the arsenal
  slot at a time, so this adds O(|Drawn|) new choices per partition, not
  a full cross-product.
- Value: arsenaled drawn card credits its future-turn tempo (a plausible
  `DrawValue`-sized number, since arsenal = "one card of known tempo
  carrying forward"). The exact constant lives in `card/effect_values.go`.
- Conflict: if the pre-draw partition already assigns a hand card to
  Arsenal, we keep the one with the higher future value.

### Phase 2 — PLAY

A drawn card that is itself an attack (or has a same-turn usable
effect) may be played this turn if:

- It's free (cost 0), OR
- We have a leftover pitch budget from earlier in the chain.

And if the chain has Go again available at the fire point (or the drawn
card itself grants Go again for later cards).

Adds a new chain-tail enumeration step: after the pre-draw chain
resolves, try appending the drawn card and each viable successor
ordering. Bounded by `|Drawn|` (at most 1 in practice).

### Phase 3 — PITCH

"Hopeful partitions": the pre-draw partition enumerator may commit an
under-funded attacker, betting that a mid-turn draw produces enough
additional pitch to cover the gap.

- Legality check is deferred to the chain leaf, *after* the draw has
  fired. A hopeful partition with an unclosable gap is pruned.
- Bound the blow-up: at most one hopeful attacker per partition, gap
  ≤ 3 (max possible pitch from a single drawn card).
- The user's causality rule is preserved: we don't peek the drawn card
  at partition time; we commit, then resolve.

### Phase 4 — DEFEND

The drawn card sits in hand until opponent's turn, where it can block
or be used as a Defense Reaction.

- Drawn card's `Defense()` + prevented-share accounting matches the
  existing defender pipeline.
- This phase swallows the HELD default: a drawn card that isn't chosen
  for defence falls back to HELD (phase 0) or ARSENAL (phase 1).

## Invariants the whole design preserves

- **No retroactive knowledge.** Pre-draw decisions (partition, chain
  order, pitch assignment) never depend on the identity of
  mid-turn-drawn cards. The solver may enumerate *hopeful* partitions
  whose legality depends on a later draw, but it doesn't read the card
  until the chain fires the draw.
- **Shared helper.** All mid-turn draw events funnel through
  `state.DrawOne()`. No per-card bookkeeping.
- **End-of-turn-draw accounting is intrinsic to HELD.** A HELD drawn
  card costs one end-of-turn refill draw. No side-channel bookkeeping
  needed — it falls out of the valuation directly (net-zero credit
  captures it).
- **Regression coverage per phase.** Each phase ships with an end-to-end
  hand test that exercises at least one representative card and pins
  the new-disposition value; earlier phases' tests stay unchanged.

## Non-goals for the first round

- Multi-draw-in-one-turn interaction (Snatch + Drawn to the Dark
  Dimension in the same chain). Phase 0 infrastructure supports it
  (the helper and slice handle N draws), but disposition logic stays
  single-draw until a later phase.
- Sigil of the Arknight's cross-turn additive draw. It keeps its
  current flat-credit path until we model cross-turn additive
  explicitly.
- Heroes with Intelligence ≠ 4. Only Viserai is implemented today;
  end-of-turn-draw accounting is hard-coded to 4 until a second hero
  lands.

## Open questions (carry through the phases)

1. `TurnSummary` exposure — do display callers want to see
   `Drawn`/`Arsenal` picks per hand, or is logging the raw numbers
   enough? (Phase 1 forces this question.)
2. Memoization — mid-turn-drawn cards shift the effective hand. Does
   `memoKey` need to track drawn cards too? First pass: mark Snatch /
   Drawn to the Dark Dimension with `NoMemo` until we've benched the
   cost of extending the key.
3. "Future-turn floor" for ARSENAL — is `DrawValue` the right number
   here, or do we want a new `ArsenalValue` constant? (Settle in
   Phase 1.)

## Scratch

*Running list of half-baked ideas, invariants I find while coding, etc.
Kept in this doc so the final PR sees just the tested outcomes.*

- (none yet)
