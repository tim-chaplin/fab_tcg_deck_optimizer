# TODO

## Simplifying assumptions that may overvalue conditional cards

Many card implementations assume optimistic outcomes for conditional effects so the simulator
doesn't have to model prior-turn state or opponent choices. These assumptions systematically
favour cards whose printed stats are backed up by conditional riders. In the future we may want
to apply a discount (e.g. 50%) to the value contributed by these conditions rather than treating
them as fully active.

- **Runechant cost discounts always apply in full.** Cards that "cost {r} less to play for each
  Runechant you control" (Reduce to Runechant, Amplify the Arknight, Drawn to the Dark Dimension,
  Rune Flash) are treated as cost 0.
- **On-hit effects always fire.** Any rider that reads "when this hits [a hero], ..." is assumed
  to always trigger (e.g. Consuming Volition's discard, Meat and Greet's Runechant, Weeping
  Battleground's arcane damage). Blocks / Dominate / attack-defence interactions aren't modelled.
- **"If you've dealt arcane damage this turn" is always true.** We assume a Runechant was in play
  before the attack and the opponent chose not to block its trigger, so the condition is
  satisfied (affects Arcanic Spike +2{p}, Consuming Volition discard rider, Meat and Greet go
  again, Sigil of Suffering +1{d}).
