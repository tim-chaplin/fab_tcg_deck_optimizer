# TODO

Running list of work we want to do on this project. Sectioned by theme.

## Simplifying Assumptions

Many card implementations assume optimistic outcomes for conditional effects so the simulator
doesn't have to model prior-turn state or opponent choices. These assumptions systematically
favour cards whose printed stats are backed up by conditional riders. In the future we may want
to apply a discount (e.g. 50%) to the value contributed by these conditions rather than treating
them as fully active.

### On-hit / combat interactions

- **Attacks always hit.** Any "when this hits a hero…" rider is assumed to fire — the sim doesn't
  model blocks. Affects Consuming Volition, Meat and Greet (on-hit Runechant), Weeping
  Battleground, Reek of Corruption, Mauvrion Skies, Runic Reaping, Oath of the Arknight. Also
  lets on-hit riders on Jack Be Nimble / Jack Be Quick (steal), Snatch (draw), Rifting (instant
  cast), Life for a Life (1{h}), Blow for a Blow (1 damage), Fervent Forerunner (Opt 2),
  Regain Composure (unfreeze) count as always firing (but the riders themselves aren't wired
  in yet — see below).
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

### "If you've dealt arcane damage this turn"

Consuming Volition, Arcanic Spike, and Meat and Greet now gate their arcane-damage riders on
`TurnState.ArcaneDamageDealt`, which playSequence flips on when a Runechant fires on an
attack/weapon and which direct-arcane cards set themselves in Play. The same clause still fires
unconditionally for:

- **Sigil of Suffering** — +1{d} on defense reactions treated as always active.

### Hero health and life-total riders

Hero health isn't tracked, so every life-gain and life-comparison rider collapses to one side.

- **Life-total comparisons never fire.** Adrenaline Rush +3{p} (less life), Down But Not Out
  (health / equipment / tokens gating), Life for a Life go-again, Blow for a Blow go-again,
  Scar for a Scar go-again, and Wounded Bull +1{p} all default to off.
- **Life-gain effects are dropped.** Healing Balm (gain 3{h}), Fyendal's Fighting Spirit,
  Sun Kiss, Sirens of Safe Harbor (graveyard 1{h}), Sigil of Fyendal (1{h} on leave),
  Fiddler's Green (3{h} entering graveyard), and the Clearwater / Restvine / Sapwood Elixir
  trio (Bloodrot Pox / Inertia / Frailty health-gain riders) all ignore the gain step.

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
- **Cross-turn aura lifecycles are collapsed.** Blessing of Occult, Sigil of Deadwood, and Sigil
  of the Arknight credit their benefits immediately (via DelayRunechants or flat damage) rather
  than modelling the full enter/leave sequence across turns. End-phase destruction clauses on
  Enchanting Melody, Sigil of Cycles, and Sigil of Fyendal are similarly dropped.
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
- **Draw / hand cycling is flattened.** Drawn to the Dark Dimension credits +3 for its draw;
  Sutcliffe's Research Notes ignores its re-ordering clause; Sink Below drops its cycling rider;
  Rise Above's alternative hand-as-cost option isn't simulated. The Emissary of Moon / Tides /
  Wind trio, Sift, Scour the Battlescape, Whisper of the Oracle (Opt), and Strategic Planning
  all similarly drop their draw / cycle steps. Trade In's discard-to-draw is dropped.
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
  Universal keyword, Test of Strength's Clash with the attacking hero, and Smashing Good Time's
  item-destruction rider all collapse to base stats plus whatever bonus we credit unconditionally.

### Assumptions that could now be dropped

Each of these predates the sim's current TurnState plumbing. The accurate check is already
available as a one-line condition; the fix is a single Play function plus a mirrored test. The
pattern to follow is the Consuming Volition fix (PR #41) for the arcane-damage items, and the
Hit the High Notes / Shrill of Skullform pattern for the aura one.

- **Sigil of Suffering** (`internal/card/runeblade/sigil_of_suffering.go`) — +1{d} buff on
  defense reactions is assumed always-true. Unlike the +{p} / on-hit riders, this one gates on
  `Defense()` which is consumed by the solver's partition scoring before `Play()` runs. Dropping
  the assumption needs a new `ConditionalDefense(state)` hook (or equivalent) so the block
  capacity can react to `state.Runechants > 0`. Yinti Yanti's defending-side +1{d} would ride
  on the same hook.
