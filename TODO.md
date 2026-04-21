# TODO

Running list of work we want to do on this project. Sectioned by theme.

## Simplifying Assumptions

Many card implementations assume optimistic outcomes for conditional effects so the simulator
doesn't have to model prior-turn state or opponent choices. These assumptions systematically
favour cards whose printed stats are backed up by conditional riders. In the future we may want
to apply a discount (e.g. 50%) to the value contributed by these conditions rather than treating
them as fully active.

### Fully model effects where we currently just credit an integer value

`internal/card/effect_values.go` centralises the damage-equivalents we use as stand-ins for
"force opponent discard" (3) and "create a Gold token" (1). These are simplifications — the
sim never actually forces a discard or tracks Gold. When we model the real state (graveyard,
Gold-token pool, opposing hand size) the rider implementations can cash out into actual
future-turn tempo instead of a flat integer, and the `effect_values.go` constants should
disappear. Mid-turn `"draw a card"` riders route through `TurnState.DrawOne` instead — see
`internal/card/card.go`. Start-of-next-turn reveal-and-put-into-hand effects (Sigil of the
Arknight) route through `card.DelayedPlay`'s `ToHand` return and land in the actual turn-2
hand rather than as a flat credit.

### LikelyToHit breadcrumbs — on-hit riders awaiting modelling

Each of the generic attacks below has a `func <card>Damage(attack int) int` helper with an
`if card.LikelyToHit(attack) { /* TODO */ }` block. The body is a placeholder so grep for
`LikelyToHit` turns these up when we come back to wire the riders. Plug the rider's
damage-equivalent into the body and remove the TODO.

- **Jack Be Quick** — on-hit steal ally (hero-specific).
- **Jack Be Nimble** — on-hit steal item (hero-specific).
- **Wreck Havoc** — on-hit DR lockout + arsenal manipulation (hero-specific).
- **Walk the Plank** — on-hit freeze target (Pirate hero-specific).
- **Tongue Tied** — on-hit arsenal face-up + banish instant (hero-specific).
- **Smash Up** — on-hit arsenal face-up + banish attack action (hero-specific).
- **Pursue to the Edge of Oblivion** — on-hit mark (hero-specific).
- **Pursue to the Pits of Despair** — on-hit mark (hero-specific).
- **Money or Your Life?** — on-hit deal 2 unless Gold given (hero-specific, Thief-repeat).
- **Humble** — on-hit hero-ability suppression.
- **Hand Behind the Pen** — on-hit arsenal face-up + banish non-attack action.
- **Fact-Finding Mission** — on-hit peek arsenal / equipment.
- **Destructive Deliberation** — on-hit create Ponder token.
- **Down But Not Out** — on-hit create Agility/Might/Vigor tokens (gated on life/equipment/token
  comparison).
- **Cut Down to Size** — on-hit conditional discard (4+ cards in hand).
- **Crash Down the Gates** — on-hit destroy top of their deck.
- **Blanch** — on-hit opponent's cards lose all colors.

Most of these require state the sim doesn't currently track (arsenal, marks, Gold / Ponder /
status tokens, opponent hand size, life totals, deck top). The relevant state-tracking gaps are
called out in the sections below — landing any of them unlocks a subset of the breadcrumbs.

### On-hit / combat interactions

- **Attacks always hit.** Any "when this hits a hero…" rider is assumed to fire — the sim doesn't
  model blocks. Affects Weeping Battleground, Reek of Corruption. Also lets on-hit riders on
  Jack Be Nimble / Jack Be Quick (steal), Snatch (draw), Rifting (instant cast), Life for a
  Life (1{h}), Blow for a Blow (1 damage), Fervent Forerunner (Opt 2), Regain Composure
  (unfreeze) count as always firing (but the riders themselves aren't wired in yet — see
  below). Mauvrion Skies, Meat and Greet, Consuming Volition, Oath of the Arknight, and
  Runic Reaping all gate their "if hits" clauses via `card.LikelyToHit` on the attack's
  printed power.
- **Dominate isn't modelled.** Drowning Dire, Overload, Pound for Pound, and Demolition Crew
  carry the keyword but the solver doesn't route around block partitioning.
- **On-hit go-again isn't granted.** Overload's on-hit clause never fires.
- **Incoming damage doesn't interrupt our auras.** Bloodspill Invocation assumes the attack lands
  before its aura gets destroyed by incoming damage, so all N Runechants are created up front.
- **Lay Low's marked-defender cost is ignored.** Card is treated as always legal and the attacker
  debuff is dropped.

### Defender-side interactions

- **"Defended by X" riders never fire.** The solver doesn't expose what card is blocking, so
  defended-by-action-card buffs/debuffs (Feisty Locals +2{p}, Freewheeling Renegades −2{p}),
  defended-by-<2-non-equipment conditions (Barraging Brawnhide, Stony Woottonhog), Out Muscle's
  equal-or-greater-power gate, and Surging Militia's +N{p} all default to off.
- **Defender-side power and defence reductions aren't modelled.** Drag Down's −3{p} attacker
  debuff and Right Behind You's defend-together +1{d} aren't simulated — the solver doesn't
  expose either side's opposing attack.
- **Defence-time instant activations are dropped.** Rally the Coast Guard and Rally the Rearguard
  carry only their printed defence; Wreck Havoc's defence-reaction lockout isn't modelled.
- **Defence-prevention and damage-prevention riders are dropped.** Battlefront Bastion's
  prevent-defence clause, Sigil of Protection's Ward N, and Enchanting Melody's incoming-damage
  prevention trigger all return only their printed stats.
- **Yinti Yanti's defending-side +1{d} is ignored.** Defence is consumed by partition scoring
  before `Play()` runs (same shape as the Sigil of Suffering hook needed below).

### Hero health and life-total riders

Hero health isn't tracked, so every life-gain and life-comparison rider collapses to one side.

- **Life-total comparisons are modelled per-hero.** Adrenaline Rush (+3{p}), Fyendal's Fighting
  Spirit (1{h} gain on attack), Life for a Life go-again, Blow for a Blow go-again, Scar for a
  Scar go-again, and Wounded Bull (+1{p}) fire when the current hero implements
  `card.LowerHealthWanter` (see `internal/card/card.go`) and stay off otherwise — a coarse proxy
  that skips per-turn life tracking. No hero opts in yet. Down But Not Out's health / equipment /
  tokens gate is not yet covered and still defaults off.
- **Life-gain effects are credited 1-to-1 with damage** when the trigger fires unconditionally
  (Healing Balm, Sun Kiss, Sirens of Safe Harbor's graveyard gain on attack, Sigil of Fyendal's
  gain on leave, Fiddler's Green's gain on defense). The Clearwater / Restvine / Sapwood Elixir
  trio (Bloodrot Pox / Inertia / Frailty riders, gated on status tokens we don't track) still
  defaults off.

### Token economies and resource trackers

- **Gold tokens aren't tracked.** Cash In, Money or Your Life, Money Where Ya Mouth Is (Wager),
  Performance Bonus, Ransack and Raze (X-cost → 0), Starting Stake, Strike Gold, Test of Strength
  (winner rider), and Wage Gold all drop their token economies.
- **Silver and Copper tokens aren't tracked.** Cash In, High Striker (Copper), and Pick a Card,
  Any Card (Silver + opponent hand inspection) default to their base effects only.
- **Action-point tracking isn't modelled.** Lead the Charge and Back Alley Breakline's
  face-up-from-deck grant are dropped.
- **Status tokens aren't created or tracked.** Infectious Host doesn't emit Frailty / Inertia /
  Bloodrot Pox; Destructive Deliberation doesn't emit Ponder; Flock of the Feather Walkers
  doesn't emit Quicken. The Elixir cycle consumes these tokens elsewhere, so their health-gain
  riders (covered above) never fire either.

### Marks, opponent state, and opponent-visible info

- **Marks aren't tracked.** Outed's +1{p} on marked heroes, Pursue to the Edge of Oblivion /
  Pursue to the Pits of Despair "mark on hit", Public Bounty's mark rider (currently credited
  unconditionally), and Relentless Pursuit's mark-plus-"attacked them this turn" gate all
  default off (or fire unconditionally where noted).
- **Opponent hand / arsenal inspection isn't modelled.** Fact-Finding Mission, Pick a Card,
  Frontline Scout's hand-peek, and Crash Down the Gates' reveal comparison drop the peek step.
- **Opponent debuffs aren't modelled.** Blanch (lose all colors), Cut Down to Size (discard),
  Humble (hero-ability suppression), and Walk the Plank (Pirate target-freezing) drop the
  debuff; the printed attack is kept.
- **Banished-zone tracking isn't modelled.** Tremor of íArathael's +2{p} rider never fires.

### Pay-extra riders and modal choices

- **"Pay {r} or lose 1{p}" is resolved as "keep power".** Bluster Buff, Chest Puff, and Look Tuff
  assume the player can always afford the upkeep.
- **Pay-to-buff-power modes are dropped.** Flex (+2{p} for {r}{r}), Punch Above Your Weight
  (+5{p} for {r}{r}{r}), and Brothers in Arms (pay-to-buff-defence) all return base stats only.
- **Modal choices are hard-coded to one branch.** Captain's Call always takes +2{p} (dropping
  the alternative Go again mode); Life of the Party's Crazy Brew substitute + random-mode
  selection isn't modelled.

### Freeze / tap state

- **Freeze / unfreeze isn't modelled.** Tit for Tat (tap/untap), Regain Composure's on-hit
  unfreeze rider, and Walk the Plank's target-freeze all default off. Tip-Off's instant discard
  activation falls in the same bucket.

### Aura / graveyard state

- **Aura presence in graveyard is assumed.** Sigil of Silphidae, Weeping Battleground, and Runic
  Fellingsong (partly) assume there's always an aura available to banish.
- **"Played or created an aura this turn" is assumed true** for cards that don't yet use the
  live `AuraCreated` / `HasPlayedType(TypeAura)` check. Reek of Corruption, Hit the High Notes,
  and Shrill of Skullform all gate this correctly; Yinti Yanti does as well. Nothing else on the
  roster currently reads the clause.
- **Cross-turn aura lifecycles are partially modelled.** `card.DelayedPlay` threads a
  PlayNextTurn callback through the deck loop for cards whose effect fires at the start of the
  owner's next action phase — Sigil of the Arknight peeks the actual post-draw top card next
  turn, and Sigil of Fyendal credits its 1{h} gain on leave the turn the aura resolves. Other
  cross-turn auras still collapse their effects into the immediate Play: Blessing of Occult
  (DelayRunechants), Sigil of Deadwood, Sigil of Silphidae (enter + leave both credited at
  play), Enchanting Melody (end-phase destruction clause dropped), Sigil of Cycles (on-leave
  discard/draw dropped).
- **Graveyard-banish additional costs are ignored.** Gravekeeping, Jack Be Nimble, Jack Be Quick,
  Looking for a Scrap, and Nimble Strike treat the banish step as free and either drop the
  rider or credit it unconditionally where noted in the card.
- **Graveyard-reorder effects aren't modelled.** Cadaverous Contraband (graveyard → top of deck)
  and Drone of Brutality (graveyard-replacement-to-deck) both drop the reorder step.

### Arsenal and hand-state effects

- **Arsenal isn't modelled.** See the dedicated Rules modelling item above — all arsenal-gated
  riders (Unmovable +1{d}, Springboard Somersault +2{d}, and the ~14 Silver Age generics listed
  there) default off or fire unconditionally where noted.
- **"No cards in hand" riders never fire.** Spring Load's +3{p} rider defaults off.
- **Draw / hand cycling is flattened.** Mid-turn draws (Snatch, Drawn to the Dark Dimension)
  route through `TurnState.DrawOne`; the drawn card competes with Held hand cards for the
  end-of-turn arsenal slot, can fund a cost shortfall via PITCH, and attaches to the chain
  tail as a free/affordable PLAY when Go again is available. Lookahead grants that scan
  `CardsRemaining` at play time (Flying High's next-attack grant, Mauvrion Skies,
  Oath of the Arknight, Runic Reaping, Condemn to Slaughter, Captain's Call) silently fizzle
  when their intended target is only drawn later in the chain — a conservative under-count
  we tolerate to avoid a retroactive re-resolution pass. Sutcliffe's Research Notes ignores its
  re-ordering clause; Sink Below drops its cycling rider; Rise Above's alternative hand-as-cost
  option isn't simulated. The Emissary of Moon / Tides / Wind trio, Sift, Scour the
  Battlescape, Whisper of the Oracle (Opt), and Strategic Planning all similarly drop their
  draw / cycle steps. Trade In's discard-to-draw is dropped.
- **Hand-on-top / hand-as-cost alternative costs aren't modelled.** Moon Wish (hand-on-top + Sun
  Kiss search) and Seek Horizon (hand-on-top + conditional go-again) pick the base mode only.
- **Fate Foreseen's "opt 1" is dropped** — block value is the printed defence only.

### Deck search, reveal, and reorder

- **Deck search isn't modelled.** Belittle (Minnowism), Nimby (Nimblism), Sound the Alarm, and
  Moon Wish's Sun Kiss search all drop the tutor step.
- **Deck reveal / peek isn't modelled.** On the Horizon's deck-peek trigger, Crash Down the
  Gates' reveal comparison and top-of-deck destruction, Ravenous Rabble's −X{p} reveal rider,
  and Demolition Crew's additional reveal cost are all collapsed.
- **Put-in-deck / deck-reorder effects aren't tracked.** Sky Fire Lanterns peeks at the top card
  (reading Deck) but any reorder step is collapsed. Warmonger's Recital's bottom-of-deck rider
  is dropped; Right Behind You's deck-bottom rider is also dropped.

### Weapon / sword chain

- **The weapon chain isn't peeked for conditional riders.** Brandish (next-weapon-attack +1{p}),
  On a Knife Edge (next-sword-attack go-again), and Visit the Blacksmith (next-sword-attack
  bonuses) drop their riders because CardsRemaining only holds action cards.

### Chain history

- **In-chain history isn't readable from Play.** Push the Point's chain-history +2{p} and
  Water the Seeds' chain-bonus for a later low-power attack both drop their riders.

### Other approximations

- **Put in Context's base-power cap is ignored** — every attack is assumed to qualify.
- **Arcane damage is counted as regular damage at creation.** Runechant tokens are credited +1
  immediately on creation; subsequent firing on attacks is purely state cleanup. Leftover tokens
  that neither fire nor carry over (end-of-sim) are slightly over-credited.
- **Regurgitating Slog's riders are fully dropped** — Play returns base power with no modelling
  attempted.
- **Uncommon keyword text is dropped.** Prime the Crowd's Crowd cheers/boos keywords, Wage Gold's
  Universal keyword, and Smashing Good Time's item-destruction rider all collapse to base stats
  plus whatever bonus we credit unconditionally. (Clash is now modelled via `card.ClashValue` —
  Test of Strength is the only card that uses it today.)

### Direction tags: undervalued vs overvalued

Bullets tag every card whose implementation drops or bends a rider, based on whether that
simplification tilts the sim's score above or below what the real card delivers. Undervalued is
the priority bucket — those cards get rejected by the optimizer before they ever get a shot at
playtesting. Cards whose assumption already appears verbatim earlier in the section are still
listed here so the direction tag is co-located with the name.

#### Undervalued (sim discounts a real card's effect)

- **Back Alley Breakline (all colours)** — face-up-from-deck action-point grant dropped; action
  points not tracked, so the tempo payoff is never scored.
- **Barraging Brawnhide (all colours)** — defended-by-<2-non-equipment +1{p} never fires; the
  card's own rider is dropped while printed power is left alone.
- **Belittle (all colours)** — dropped Minnowism tutor step; the card prints go-again and searches
  a specific follow-up, and only the go-again lands.
- **Blanch (all colours)** — dropped on-hit "cards lose all colors" debuff; effectively vanilla
  power.
- **Brandish (all colours)** — dropped next-weapon-attack +1{p}; with go-again printed the rider
  would typically cash out.
- **Brothers in Arms (all colours)** — dropped pay-to-buff-defence rider; card never gets credit
  for its +2{d} swing.
- **Cadaverous Contraband** — dropped graveyard-top-of-deck fix-up; the card's whole point is
  unavailable.
- **Captain's Call (all colours)** — modal pick is hard-coded to +2{p}; the alternative go-again
  mode that could chain a bigger attack is never chosen.
- **Cash In (all colours)** — Gold / Silver / Copper economy plus draw riders dropped; card
  returns only base stats.
- **Clearwater / Restvine / Sapwood Elixir (all colours)** — health-gain rider dropped (still
  credits the +{p} scan, so the loss is only the life side).
- **Crash Down the Gates (all colours)** — on-hit deck destruction + reveal comparison dropped.
- **Cut Down to Size (all colours)** — on-hit opponent discard dropped.
- **Destructive Deliberation (all colours)** — Ponder-token creation dropped entirely.
- **Down But Not Out (all colours)** — none of Agility / Might / Vigor token branches fire.
- **Drone of Brutality (all colours)** — graveyard-replacement-to-deck rider dropped.
- **Emissary of Moon (all colours)** — hand-cycle draw dropped.
- **Emissary of Tides (all colours)** — hand-cycle-for-+2{p} dropped.
- **Emissary of Wind (all colours)** — hand-cycle-for-go-again dropped.
- **Enchanting Melody (all colours)** — damage-prevention trigger dropped (aura-created flag is
  the only value credited).
- **Fact-Finding Mission (all colours)** — opponent arsenal/equipment peek dropped.
- **Fate Foreseen (all colours)** — Opt 1 dropped; block value is printed defence only.
- **Feisty Locals (all colours)** — defended-by-action +2{p} rider never fires.
- **Fervent Forerunner (all colours)** — on-hit Opt 2 and the from-arsenal go-again both dropped.
- **Flex (all colours)** — pay-{r}{r}-for-+2{p} mode dropped; only printed power.
- **Flock of the Feather Walkers (all colours)** — Quicken token creation dropped (and the
  reveal cost too).
- **Frontline Scout (all colours)** — hand-peek plus arsenal-only go-again dropped.
- **Fyendal's Fighting Spirit (all colours)** — on-defend 1{h} gain dropped; on-attack gain is
  modelled for `card.LowerHealthWanter` heroes only.
- **Gravekeeping (all colours)** — graveyard-banish additional value dropped.
- **Hand Behind the Pen (all colours)** — on-hit arsenal / banish-instant dropped.
- **High Striker (all colours)** — Copper-token economy dropped.
- **Humble (all colours)** — hero-ability suppression debuff dropped.
- **Infectious Host (all colours)** — Frailty / Inertia / Bloodrot Pox token emission dropped
  (and the Elixir cycle that would consume them).
- **Jack Be Nimble (all colours)** — graveyard-banish +1{p} / go-again and on-hit item steal
  dropped.
- **Jack Be Quick (all colours)** — graveyard-banish +1{p} / go-again and on-hit ally steal
  dropped.
- **Lead the Charge (all colours)** — face-up-from-deck action-point grant dropped.
- **Life of the Party (all colours)** — all three modes default off including go-again; Crazy
  Brew substitute never fires.
- **Looking for a Scrap (all colours)** — graveyard-banish bonus dropped.
- **Money or Your Life (all colours)** — Gold-exchange rider dropped; repeatable mode never
  fires.
- **Money Where Ya Mouth Is (all colours)** — Wager Gold payout dropped (scan bonus is still
  credited).
- **Moon Wish (all colours)** — hand-on-top alt cost plus Sun Kiss tutor dropped.
- **Nimble Strike (all colours)** — graveyard-banish bonus dropped.
- **Nimby (all colours)** — Nimblism tutor dropped.
- **On a Knife Edge (all colours)** — next-sword-attack go-again dropped (weapon chain not
  scanned).
- **On the Horizon (all colours)** — deck-peek trigger dropped.
- **Out Muscle (all colours)** — defender-power go-again gate dropped; the printed rider never
  fires in the sim.
- **Outed (all colours)** — +1{p} vs marked hero never fires.
- **Performance Bonus (all colours)** — arsenal-conditional go-again dropped (on-hit Gold is
  credited).
- **Pick a Card, Any Card (all colours)** — opponent hand inspection and Silver-token rider
  dropped.
- **Promise of Plenty (all colours)** — arsenal-placement rider plus arsenal go-again dropped.
- **Pursue to the Edge of Oblivion (all colours)** — on-hit mark dropped.
- **Pursue to the Pits of Despair (all colours)** — on-hit mark dropped.
- **Push the Point (all colours)** — chain-history +2{p} dropped (chain history unreadable from
  Play).
- **Punch Above Your Weight (all colours)** — pay-{r}{r}{r}-for-+5{p} dropped.
- **Rally the Coast Guard (all colours)** — defence-time instant activation dropped.
- **Rally the Rearguard (all colours)** — defence-time instant activation dropped.
- **Ransack and Raze (all colours)** — Gold / Landmarks X-cost treated as 0 so the ramp pay-off
  never lands.
- **Regain Composure (all colours)** — on-hit unfreeze dropped (scan-bonus itself is credited).
- **Regurgitating Slog (all colours)** — Sloggism graveyard-banish Dominate rider fully dropped.
- **Relentless Pursuit** — marked-target + attacked-them-this-turn chain rider dropped.
- **Rifting (all colours)** — on-hit instant cast dropped.
- **Right Behind You (all colours)** — defend-together +1{d} and deck-bottom rider dropped.
- **Rise Above (all colours)** — hand-as-cost alt dropped (card fails cost check unless printed
  cost is met).
- **Runeblade Condemn to Slaughter (all colours)** — aura-trade rider and opponent-aura
  destruction dropped (only same-turn Runeblade-attack +N is modelled).
- **Runeblade Runic Fellingsong (all colours)** — cannot credit BOTH the printed 1 arcane AND
  the graveyard-banish rider; only one fires.
- **Runeblade Splintering Deadwood (all colours)** — aura-swap modelled as net-zero (no credit
  for the tempo of trading a weak aura for a Runechant).
- **Runeblade Sutcliffe's Research Notes (all colours)** — top-of-deck re-ordering clause
  dropped.
- **Scour the Battlescape (all colours)** — hand-cycle plus arsenal go-again dropped.
- **Seek Horizon (all colours)** — hand-on-top alt cost plus conditional go-again dropped.
- **Sift (all colours)** — hand cycling dropped.
- **Sigil of Cycles (all colours)** — end-phase discard-and-draw dropped.
- **Sigil of Protection (all colours)** — Ward N dropped.
- **Smash Up (all colours)** — on-hit arsenal face-up + banish attack action dropped.
- **Snatch (all colours)** — currently fires the on-hit draw via TurnState.DrawOne, but a drawn
  card only recovers part of a real draw's value in the sim (no cross-turn shuffle benefit).
- **Sound the Alarm (all colours)** — deck-search rider dropped.
- **Spring Load (all colours)** — +3{p} empty-hand rider never fires.
- **Springboard Somersault (all colours)** — arsenal-only +2{d} never fires.
- **Starting Stake (all colours)** — Gold-token economy dropped.
- **Strategic Planning (all colours)** — graveyard recovery plus end-phase draw dropped.
- **Stony Woottonhog (all colours)** — defended-by-<2-non-equipment rider dropped (same shape as
  Barraging Brawnhide below).
- **Sun Kiss (all colours)** — Moon Wish synergy (draw + go again) dropped; the 3{h} gain is
  modelled.
- **Surging Militia (all colours)** — defended-by +N{p} rider dropped.
- **Tip-Off (all colours)** — instant discard activation dropped.
- **Tongue Tied (all colours)** — on-hit arsenal face-up + banish instant dropped.
- **Trade In (all colours)** — discard-to-draw plus arsenal-only go-again dropped.
- **Tremor of íArathael (all colours)** — banished-zone +2{p} never fires.
- **Unmovable (all colours)** — arsenal +1{d} never fires.
- **Visit the Blacksmith (all colours)** — next-sword-attack bonuses dropped.
- **Wage Gold (all colours)** — Universal keyword plus Gold wager dropped.
- **Walk the Plank (all colours)** — Pirate-specific target-freeze dropped.
- **Warmonger's Recital (all colours)** — bottom-of-deck rider dropped (scan bonus credited).
- **Whisper of the Oracle (all colours)** — Opt dropped.
- **Wreck Havoc (all colours)** — defence-reaction lockout plus arsenal-banish dropped.
- **Yinti Yanti (all colours)** — defending-side +1{d} dropped (aura-created clause is modelled).

#### Overvalued (sim credits more than the real card delivers)

- **Bluster Buff** — pay {r}-or-lose-1{p} resolved as "always pay"; player short on pitch would
  lose the point, so base power is over-credited.
- **Chest Puff** — pay {r}-or-lose-1{p} resolved as "always pay"; same shape as Bluster Buff.
- **Lay Low (all colours)** — treated as always legal even without a marked defender; real card
  is uncastable when no hero is marked.
- **Look Tuff** — pay {r}-or-lose-1{p} resolved as "always pay"; same shape as Bluster Buff.
- **Plunder Run (all colours)** — scan-target +N is credited unconditionally instead of only
  when the played-from-arsenal gate is met.
- **Public Bounty (all colours)** — mark rider fires unconditionally; real card requires the
  opponent be marked and the rider is tied to the mark.
- **Put in Context** — base-power cap on what it can block is ignored; every attack is assumed
  to qualify so the defence is always live.
- **Runeblade Arcanic Crackle** — printed 1 arcane added unconditionally; Ward on the opponent
  would prevent part of it.
- **Runeblade Bloodspill Invocation (all colours)** — all N Runechants credited up front;
  incoming damage this turn would destroy the aura before the payoff fires.
- **Runeblade Drowning Dire (all colours)** — Dominate keyword + AuraCreated flag both dropped
  while the printed power is left alone — partial under-count, but in contexts where a
  follow-up aura-reading card is present, Drowning Dire is materially overstated because the
  Dominate protection of that damage is not modelled.
- **Runeblade Sigil of Silphidae** — assumed we always have an aura to banish on both enter and
  leave, crediting 2 damage; real card fizzles both triggers when the graveyard is dry.
- **Runeblade Singeing Steelblade (all colours)** — printed 1 arcane added unconditionally; Ward
  would prevent it.
- **Runeblade Weeping Battleground (all colours)** — assumes the graveyard always has a
  banishable aura; fizzles when empty.
- **Scout the Periphery (all colours)** — scan bonus credited whenever a target exists, even
  though the real rider requires the card to be played from arsenal.
- **Smashing Good Time (all colours)** — scan bonus credited unconditionally (same shape as
  Scout the Periphery) rather than requiring the arsenal gate.

#### Likely neutral

- **Aether Slash (all colours)** — printed 1 arcane doubles as the text-rider damage; the code
  adds it once when the non-attack pitch condition is met, which mirrors the card's actual
  outcome on average.
- **Blessing of Occult (all colours)** — tokens routed to next-turn carryover; tempo value is
  preserved, just shifted a turn.
- **Come to Fight / Minnowism / Nimblism / Sloggism (all colours)** — scan next attack and
  credit +N only when a target actually exists; matches how the card plays in practice when the
  chain is well-ordered.
- **Deathly Duet (all colours)** — Pitched scan may let both riders fire even without the
  specific attribution; on average this approximates the real outcome when the hero pitches one
  of each type.
- **Drawn to the Dark Dimension (all colours)** — draw rider routed through TurnState.DrawOne
  and variable cost respected; accurate modelling rather than a simplification.
- **Flying High (all colours)** — +1{p} matching-colour rider is modelled via CardsRemaining
  scan plus GrantedGoAgain; both clauses fire only when a qualifying future attack exists.
- **Hit the High Notes / Shrill of Skullform / Vantage Point / Runerager Swarm (all colours)** —
  aura-created check is modelled directly via HasAuraInPlay; result matches the printed gate.
- **Hocus Pocus / Spellblade Strike / Spellblade Assault / Reduce to Runechant / Read the
  Runes** — Runechant = +1 future damage identity; faithful modelling of token creation.
- **Malefic Incantation (all colours) / Runeblood Incantation (all colours)** — split between
  carryover and in-turn credit avoids double-counting the same rune as both current-turn
  discount and future-turn damage.
- **Sigil of Deadwood (all colours)** — Runechant deliberately delayed to next-turn carryover
  rather than this turn's discount pool.
- **Sky Fire Lanterns (all colours)** — peeks actual deck top and matches on pitch; neither
  over- nor under-credits.
- **Sutcliffe's Research Notes (all colours)** — also listed under Undervalued for the dropped
  re-order clause; the reveal-and-count core is accurate.
- **Trot Along / Water the Seeds (all colours)** — scan CardsRemaining for a qualifying target
  and fizzle when none exists; mirrors real card's interaction with chain order.

Cards not listed in any bucket (Critical Strike, Brutal Assault, Raging Onslaught, Muscle Mutt,
Wounding Blow, Dodge, Evasive Leap, Toughen Up, Rune Flash, Amplify the Arknight, Fragile Aura
helper, Aura Helper, Next-Attack-Action helper) have no printed rider beyond the base stat
line, so there's nothing to tag — they score exactly what they print.
