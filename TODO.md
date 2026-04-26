# TODO

Running list of work we want to do on this project. Sectioned by theme.

Per-card unimplemented riders are now annotated directly on the card files via the
`card.NotImplemented` marker plus a `// not implemented: <quirk>` comment above it. To
audit what's still rough on a given card, open the file. The sections below describe the
broader state-tracking and framework-level gaps that gate multiple cards at once —
landing any of them lets several `NotImplemented` markers come off in one pass.

### Damage-equivalent constants in `effect_values.go`

`internal/card/effect_values.go` centralises the damage-equivalents we use as stand-ins for
"force opponent discard" (3) and "create a Gold token" (0). These are simplifications — the
sim never actually forces a discard or tracks Gold. When we model the real state (graveyard,
Gold-token pool, opposing hand size) the rider implementations can cash out into actual
future-turn draw instead of a flat integer, and the `effect_values.go` constants should
disappear.

### State-tracking gaps that gate multiple cards

These are the systemic features the sim doesn't model yet. Each gates a bucket of `// not
implemented` riders across the card roster.

- **Hero health and life-total tracking.** No per-turn hero-life accounting. Life-comparison
  riders ("if you have less {h} than an opposing hero") use the `card.LowerHealthWanter`
  hero-attribute proxy: the rider fires for heroes that opt in and never fires otherwise.
  Life-gain effects are credited 1-to-1 with damage at trigger time. Modelling real life
  totals would let conditional grants fire correctly per-turn instead of per-hero.
- **Gold / Silver / Copper / Quicken / Ponder / Frailty / Inertia / Bloodrot Pox token
  economies.** None are tracked. Cards that mint or consume these tokens collapse to base
  stats or a flat damage-equivalent. Adding a per-token counter on `TurnState` plus a
  destroy-and-redeem hook would unblock the bulk of the affected riders.
- **Action-point tracking.** The sim doesn't track action points; cards that grant them
  drop the tempo payoff entirely.
- **Marks and "attacked them this turn" tracking.** No per-hero mark state. Cards that gate
  on a marked defender / attacker fall back to credit-unconditionally or drop-unconditionally.
- **Opponent hand / arsenal / banished-zone visibility.** The sim doesn't expose the
  opposing player's hand, arsenal, or banished zone, so peek / inspection / count riders
  collapse.
- **Freeze and tap state.** No tap/untap counter; freeze and unfreeze riders default off.
- **Defender-side hooks during attacks.** The solver consumes `Defense()` before `Play()`
  runs and doesn't expose what card is blocking nor reduce the attacker's power
  defender-side. Riders keyed on "defended by X", "defended by < N non-equipment", or
  defender-side debuffs need a defender-aware Play hook to land.
- **Defence-prevention and damage-prevention triggers.** No "prevent the next N damage"
  state; cards that grant Ward N or pre-emptive prevention return only their printed stats.
- **Defence-time instant activations.** Cards whose printed text adds a chain-link defender
  or activates an instant during an attack chain carry only their printed defence today.
- **Pay-extra / modal cost choices.** "Pay {r} or lose 1{p}", "pay {r}{r} for +N{p}",
  "choose go-again or +N{p}", and Crazy Brew substitutes don't probe the resource budget;
  they pick one branch and stick with it. A pay-aware modal cost evaluator would let the
  solver pick the best mode per partition.
- **Hand-on-top / hand-as-cost alternative costs.** "Put a card on top of your deck rather
  than pay {r}" isn't modelled — cards fall back to their printed cost.
- **Mid-turn draw side-channels.** `TurnState.DrawOne` puts drawn cards into `Drawn` for
  carry-as-Held or arsenal promotion, but drawn cards can't pitch or extend the attack
  chain (would leak top-of-deck identity into the solver's line choice). Lookahead grants
  that scan `CardsRemaining` silently fizzle when their target is drawn rather than in the
  starting hand — a conservative under-count we tolerate.
- **Graveyard-banish additional costs.** Several cards have "as an additional cost,
  banish a card from your graveyard" riders that the sim treats as free — the banish step
  isn't evaluated against actual graveyard contents.
- **Graveyard-reorder and put-on-top-of-deck effects.** No deck-top mutation pipe.
- **Deck-search tutors.** Belittle's Minnowism, Nimby's Nimblism, Sound the Alarm, Moon
  Wish's Sun Kiss search — the tutor step drops, even when the searched card is in deck.
- **Top-of-deck reveal and reorder.** Some cards peek `s.Deck` (Sky Fire Lanterns,
  Ravenous Rabble, On the Horizon) but reorder steps are collapsed; reveal-comparison
  riders like Crash Down the Gates collapse too.
- **Weapon chain visibility from `Play`.** `CardsRemaining` only carries action cards;
  weapon swings aren't visible to look-ahead riders that gate on "next sword attack" /
  "next weapon attack". Brandish, On a Knife Edge, Visit the Blacksmith all drop their
  riders.
- **In-chain history readable from Play.** A card's `Play` doesn't see what played
  earlier in this same chain (it sees `CardsPlayed` from earlier resolutions but not
  immediate-prior chain history needed for chain-history riders like Push the Point and
  Water the Seeds).
- **Aura-created vs aura-played semantics.** `TurnState.HasAuraInPlay` covers most "have
  you played or created an aura this turn" reads, but a few specialised aura-state
  questions (e.g. trade-an-aura-for-a-runechant value) aren't surfaced.
- **Arcane damage credited on Runechant creation.** Runechant tokens are credited +1
  damage-equivalent at creation rather than on fire; leftover tokens at end-of-sim are
  slightly over-credited (rare in practice).

### LikelyToHit / EffectiveAttack notes

- `EffectiveAttack` (printed `Card.Attack()` + `BonusAttack`, clamped at 0) is the canonical
  attack-power read for hit-likelihood checks. `LikelyToHit(self)` folds it in along with
  `EffectiveDominate`. Granters set `pc.BonusAttack += N` on the target's `CardState`
  rather than returning the bonus from their own `Play` — the +N attributes to the buffed
  attack's chain slot, and any "if this hits" rider on the target reads the buffed value.
- For grants whose "if this hits" rider needs to see the target's *fully-resolved* attack
  state (post-grants from later cards in the chain), use `AddEphemeralAttackTrigger` —
  Mauvrion Skies and Runic Reaping route their on-hit Runechant clauses this way.
