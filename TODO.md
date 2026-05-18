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
- **Gold / Silver / Copper / Quicken / Frailty / Inertia / Bloodrot Pox token economies.**
  None are tracked. Cards that mint or consume these tokens collapse to base stats or a
  flat damage-equivalent. Each new token kind needs a `TokenType` enum entry in
  `internal/sim/tokens.go`, a handler describing the destroy condition, and a `s.CreateX`
  helper paired with the relevant card-mint sites.
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
  or activates an instant during an attack chain carry only their printed defence.
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
- **Aura-created vs aura-played semantics.** `TurnState.HasPlayedOrCreatedAura` covers most "have
  you played or created an aura this turn" reads, but a few specialised aura-state
  questions (e.g. trade-an-aura-for-a-runechant value) aren't surfaced.
- **Arcane damage credited on Runechant creation.** Runechant tokens are credited +1
  damage-equivalent at creation rather than on fire; leftover tokens at end-of-sim are
  slightly over-credited (rare in practice).

### Unimplemented cards by feature

The 64 cards left in `internal/cards/notimplemented/` grouped by the systemic feature
that gates them. Counts are by file (one card-name; each file typically holds 1–3
pitch-color variants). Landing one bullet typically lets every card in the bucket come
out of `notimplemented/` together.

- **Arsenal / item permanent manipulation** (9): on-hit opponent-arsenal poke,
  arsenal-fill-from-deck-top end-phase, arsenal-wipe on hit, on-hit item destruction,
  defense-reaction lockout. Needs a write-side "destroy / replace permanents in play"
  hook beyond what `DestroyAura` covers, plus an opposing-arsenal model.
- **Status tokens** (9): Frailty, Inertia, Bloodrot Pox, freeze / unfreeze, Crowd
  cheers / boos, passive-tap pirates. Each is a stateful counter on a hero / card that
  riders read or destroy. Mirrors the existing token plumbing (Runechant, Gold, …) but
  on a per-hero / per-card axis instead of a global pool.
- **Hand cycling** (6): "discard a card from hand, draw a card" mid-chain, with riders
  that fire conditional on the discard happening. Needs a `DiscardFromHand` primitive
  and a way to gate "if you cycled a card" riders.
- **Defense-time activated abilities / DR refinements** (5): defense-time instants
  (Rally the Coast Guard / Rearguard), Instant +N{d} grants to a defending attack
  action card, base-power caps on what a Block can defend, Instant arcane-prevention
  scaled by deck reveal. Defenders today resolve as one-shot Plays inside
  `defendersDamage`; landing this would give the defense phase its own chain.
- **Reactive talismans** (4): Spellvoid passive, AR-buff-event react, opponent-draw-event
  react, pitch-1-event react. All four self-destroy on a specific event; the sim has no
  passive event hook today, so the talismans collapse to base stats.
- **Weapon-chain inspection** (3): "next sword / weapon attack +N{p} or go again". Today
  `CardsRemaining` only carries action cards; weapon swings aren't visible to look-ahead
  riders. Needs to extend the look-ahead view past the action portion of the chain.
- **Opponent state observation** (3): peek opponent hand / arsenal / equipment. The sim
  doesn't model the opposing player's private zones at all, so peek riders collapse.
- **Token economies** (3): Gold / Silver / Landmark mints by trigger (discard creates
  Gold, etc.). Slots into the existing token-creation plumbing; each new mint site is
  a small addition.
- **Action-point-from-deck-reveal** (2): cards that grant a free chain step gated on
  the top card of the deck. Needs a deck-top-reveal-into-AP path.
- **Aura-trade / aura-destroy-opponent** (2): "destroy an opposing aura" or "trade aura
  for aura". Opposing-aura state isn't modelled.
- **Banished-zone count** (1): "+N{p} per card in your banished zone" needs a per-game
  banish counter (not just this-turn `s.Banish`); self-destroys-on-play-from-banished
  needs a was-played-from-banished marker on `CardState`.
- **Hand-as-alt-cost** (2): "put a card on top of deck instead of paying {r}". Today the
  partition only knows the printed cost; an alt-cost evaluator would have to consider
  hand-card alt-costs alongside pitch.
- **Modal pay-extra-for-bonus** (2): pay {r}{r} for +2{p}, pay {r}{r}{r} for +5{p}. The
  existing `ModalCost` framework supports per-mode costs but the upstream cards aren't
  wired through; landing this is mostly per-card work.
- **Mark mechanic remnants** (2): Lay Low's marked-defender attacker debuff, Tip-Off's
  instant-discard-to-mark activation. Both need state we don't have (our hero being
  marked / discard-from-hand-as-cost).
- **Hero-ability suppression / global debuffs** (2): "opponent loses all colors", "hero
  ability does nothing this turn". Needs a turn-scoped opt-out flag the hero ability
  pipeline reads.
- **Quicken tokens** (2): Opt-and-quicken-on-reveal, reveal-cost-with-quicken.
- **Ward (opponent damage prevention)** (1): aura that reduces opponent's incoming
  damage to our hero. New persistent state; one card.
- **Damage-prevention triggers** (1): "when you prevent damage this turn, fire X".
- **Deck-top reveal compare + destroy** (1): Crash Down the Gates' reveal-and-cap rider.
- **Health / Overpower / agility-might-vigor tokens** (1): Down But Not Out's stack of
  tokens-and-Overpower comparisons.
- **Chain-history-readable Play** (1): Push the Point's "+2{p} if a card has been played"
  rider. `CardsPlayed` is read-after-the-fact; the rider needs a mid-chain history pipe.
- **Self-destroying sigils with leave-arena triggers** (1): Sigil of Cycles' "destroy at
  start of action phase, leaves arena → discard then draw".
- **Misc unique mechanics** (1): on-hit Wizard instant-casting grant (Rifting).

### Same-turn item activation

Items created mid-chain can't be spent the same turn — only items carried in via
`priorItems` participate in the wmask's activated-ability enumeration. So a chain like
"Strike Gold creates Gold on hit → spend Gold for {2}, draw a card" doesn't get found:
the Gold lands in `s.Items` for next turn but no chain step in the current turn fires
its ability. Items are slightly underpowered as a result. Fixing this means letting the
chain runner's activated-ability list mutate dynamically with token creates rather than
being committed at chain start.

### Weapons are Cards

The Weapon interface includes the Card interface; weapons are sometimes cardlike (they
have attack power, which can be buffed, they can be granted Go Again, etc.) but are also
different from cards (they're never played, drawn from the deck, pitched, etc.). They
should really be treated as a completely separate type. However, parts of the sim currently
treat Weapons as Cards, so that will have to be carefully disentangled.

`internal/registry/ids/weapon_ids.go` aliases `WeaponID = CardID` and anchors the weapon
constants at `FakeHugeAttack + iota + 1` so they don't collide with card / fake IDs in the
shared cache slots. Ideally weapons would have their own `WeaponID uint16` type starting at
1, separate from `CardID`. Blocked by depth: every weapon swing flows through the same
chain runner as deck cards (`bestSequence` permutes one `[]card.Card` slice; weapons rely
on `*card.CardState` for `BonusAttack` / `GrantedGoAgain` and call helpers like
`s.ApplyAndLogEffectiveAttack(self)` / `s.ApplyAndLogRiderOnPlay(self, …)` that read
`self.Card.*`; the chain step / display name / attacker meta caches are keyed by `CardID`).
Splitting the type cleanly needs either a slot-tagged permutation that branches per-step
between card and weapon paths, or a parallel `WeaponState` + parallel helpers — ~200–300
lines across `card/`, `weapon/`, `hand/` plus every weapon impl.

### LikelyToHit / EffectiveAttack notes

- `EffectiveAttack` (printed `Card.Attack()` + `BonusAttack`, clamped at 0) is the canonical
  attack-power read for hit-likelihood checks. `LikelyToHit(self)` folds it in along with
  `EffectiveDominate`. Granters set `pc.BonusAttack += N` on the target's `CardState`
  rather than returning the bonus from their own `Play` — the +N attributes to the buffed
  attack's chain slot, and any "if this hits" rider on the target reads the buffed value.
- For grants whose "if this hits" rider needs to see the target's *fully-resolved* attack
  state (post-grants from later cards in the chain), append the rider closure to the
  target's `CardState.OnHit` — Mauvrion Skies and Runic Reaping route their on-hit
  Runechant clauses this way. The chain runner fires `OnHit` post-AR-buff.

### Tech debt

- `internal/sim/exports_test.go`: drop the test-only re-exports. Migrate the remaining `internal/sim/*_test.go` cases that depend on them over to `turntests/` (which drives via `(*Deck).EvalOneTurnForTesting`) so the package stops exposing internals.
- consolidate card definitions into a single top-level folder. Today they're split across `internal/card/cards/`, `internal/hero/heroes/`, `internal/weapon/weapons/` — three siblings under three different parent packages.
- move `internal/card/types.go` (CardType / TypeSet bitfield) to `internal/sim/card_types.go` — it's a sim-side optimisation, not a card-interface concern.
- move `internal/testutils/` files to `*_test.go` names so they don't get compiled into the main binary.
- audit everything under `internal/sim/` and see if it makes sense where it is.
- combine `internal/sim/hand_aura_test.go` and `internal/sim/deck_aura_test.go` into `internal/sim/aura_test.go`. Mid-turn-draw is split across packages (`internal/sim/hand_mid_turn_draw_test.go` + `turntests/deck_mid_turn_draw_test.go`); fold the sim-side cases into `turntests/` if possible.
- remove the local `zeroDefenseAura` fake from `turntests/weeping_battleground_test.go` once `defendersDamage` seeds the defense-phase state graveyard from `priorGraveyard + defenders` (currently defenders-only, asymmetric with the attack-phase seed). After that, the test can put a real aura in `prior.Graveyard` instead of routing one through a 0-defense plain block.
- clean stale `TurnState` references in `internal/sim/cardmeta.go` comments — the type was removed; only `TurnSummary` survives.