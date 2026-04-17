# TODO

## Simplifying assumptions that may overvalue conditional cards

Many card implementations assume optimistic outcomes for conditional effects so the simulator
doesn't have to model prior-turn state or opponent choices. These assumptions systematically
favour cards whose printed stats are backed up by conditional riders. In the future we may want
to apply a discount (e.g. 50%) to the value contributed by these conditions rather than treating
them as fully active.

### On-hit / combat interactions

- **Attacks always hit.** Any "when this hits a hero…" rider is assumed to fire — the sim doesn't
  model blocks. Affects Consuming Volition, Meat and Greet (on-hit Runechant), Weeping
  Battleground, Reek of Corruption, Mauvrion Skies, Runic Reaping, Oath of the Arknight.
- **Dominate isn't modelled.** Drowning Dire returns its printed attack only.
- **Incoming damage doesn't interrupt our auras.** Bloodspill Invocation assumes the attack lands
  before its aura gets destroyed by incoming damage, so all N Runechants are created up front.
- **Lay Low's marked-defender cost is ignored.** Card is treated as always legal and the attacker
  debuff is dropped.

### "If you've dealt arcane damage this turn"

Consuming Volition now gates its discard rider on `state.Runechants > 0` at Play time (see
PR #41). The same clause still fires unconditionally for:

- **Arcanic Spike** — +2{p} baked into Attack().
- **Meat and Greet** — conditional Go again treated as printed.
- **Sigil of Suffering** — +1{d} on defense reactions treated as always active.

### Aura / graveyard state

- **Aura presence in graveyard is assumed.** Sigil of Silphidae, Weeping Battleground, and Runic
  Fellingsong (partly) assume there's always an aura available to banish.
- **"Played or created an aura this turn" is assumed true** for cards that don't yet use the
  live `AuraCreated` / `HasPlayedType(TypeAura)` check. Reek of Corruption, Hit the High Notes,
  and Shrill of Skullform all gate this correctly; nothing else on the roster currently reads the
  clause.
- **Cross-turn aura lifecycles are collapsed.** Blessing of Occult, Sigil of Deadwood, and Sigil
  of the Arknight credit their benefits immediately (via DelayRunechants or flat damage) rather
  than modelling the full enter/leave sequence across turns.

### Arsenal / hand-state effects

- **Arsenal isn't modelled.** Unmovable and Springboard Somersault never trigger their arsenal
  riders (+1{d} / +2{d}). Springboard's cost assumes hand-play.
- **Draw / hand cycling is flattened.** Drawn to the Dark Dimension credits +3 for its draw;
  Sutcliffe's Research Notes ignores its re-ordering clause; Sink Below drops its cycling rider;
  Rise Above's alternative hand-as-cost option isn't simulated.
- **Fate Foreseen's "opt 1" is dropped** — block value is the printed defence only.

### Other approximations

- **Put in Context's base-power cap is ignored** — every attack is assumed to qualify.
- **Arcane damage is counted as regular damage at creation.** Runechant tokens are credited +1
  immediately on creation; subsequent firing on attacks is purely state cleanup. Leftover tokens
  that neither fire nor carry over (end-of-sim) are slightly over-credited.
- **Put-in-deck / deck-reorder effects aren't tracked.** Sky Fire Lanterns peeks at the top card
  (reading Deck) but any reorder step is collapsed.

## Assumptions that could now be dropped

Each of these predates the sim's current TurnState plumbing. The accurate check is already
available as a one-line condition; the fix is a single Play function plus a mirrored test. The
pattern to follow is the Consuming Volition fix (PR #41) for the arcane-damage items, and the
Hit the High Notes / Shrill of Skullform pattern for the aura one.

- **Arcanic Spike** (`internal/card/runeblade/arcanic_spike.go`) — +2{p} rider for "dealt arcane
  damage this turn" is baked into Attack(). Gate on `state.Runechants > 0` at Play time; return
  the un-bonused attack otherwise.
- **Meat and Greet** (`internal/card/runeblade/meat_and_greet.go`) — "dealt arcane damage" Go
  again clause is assumed always-true. Gate the Go again on `state.Runechants > 0` via
  `PlayedCard.GrantedGoAgain`.
- **Sigil of Suffering** (`internal/card/runeblade/sigil_of_suffering.go`) — +1{d} buff on
  defense reactions is assumed always-true. Unlike the +{p} / on-hit riders, this one gates on
  `Defense()` which is consumed by the solver's partition scoring before `Play()` runs. Dropping
  the assumption needs a new `ConditionalDefense(state)` hook (or equivalent) so the block
  capacity can react to `state.Runechants > 0`.
