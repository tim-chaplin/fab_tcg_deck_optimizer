package sim

// Test-only entry point: drives two consecutive turns against a fixed deck order so
// tests can assert on cross-turn effects (start-of-turn-aura fires, banished /
// graveyard persistence, recycle into deck) without running the multi-shuffle Evaluate
// loop. Cross-turn state propagation goes through the same per-turn helpers the deck
// loop uses, so asserted behaviour matches what the full simulator would produce.

// TwoTurnSummary is the result of EvalTwoTurnsForTesting: TurnSummary snapshots after
// turn 1 and after turn 2. Use Turn1 for the first turn's chain outcome and what's
// queued for turn 2; use Turn2 for cross-turn effects — start-of-turn-aura damage
// folded into Turn2.Value, fires recorded in Turn2.TriggersFromLastTurn, exhausted
// auras and chain cards in Turn2.State.Graveyard, banished cards persisted across the
// boundary.
type TwoTurnSummary struct {
	Turn1 TurnSummary
	Turn2 TurnSummary
}

// EvalTwoTurnsForTesting runs two turns sequentially against d.Cards in source order
// (no shuffle). initial seeds turn 1's cross-turn state; initialHand is turn 1's
// starting hand (see EvalOneTurnForTesting for the seed semantics). Turn 2 inherits
// turn 1's outputs through the deck loop's carry/recycle path. If turn 2 can't be
// dealt (deck exhausted) the returned Turn2 is the zero TurnSummary; Turn1 is still
// populated.
func (d *Deck) EvalTwoTurnsForTesting(mp Matchup, initial TurnState, initialHand []Card) TwoTurnSummary {
	CurrentHero = d.Hero
	handSize := d.Hero.Intelligence()
	if handSize <= 0 {
		return TwoTurnSummary{}
	}

	turn1Hand, head, ok := resolveTurn1Hand(d.Cards, initialHand, handSize)
	if !ok {
		return TwoTurnSummary{}
	}

	deckSize := len(d.Cards)
	buf := make([]Card, deckSize*2+handSize*2)
	copy(buf, d.Cards)
	handBuf := make([]Card, handSize, handSize+startOfTurnRevealRoom)
	tail := deckSize

	carry := turnCarryFromInitial(initial, head, tail, handSize)

	// Turn 1: caller-supplied hand, drawCount=0.
	h := handBuf[:len(turn1Hand)]
	copy(h, turn1Hand)
	turn1, _, _ := runChainAfterDeal(&carry, buf, h, 0, d.Hero, d.Weapons, mp, nil)
	carry.applyTurnAndRotate(turn1, buf)

	// Turn 2: dealt through runOneTurnInShuffle so cross-turn wiring (recycle, held
	// carry, banished / graveyard / aura propagation) matches the deck loop.
	turn2, _, _, ok := runOneTurnInShuffle(&carry, buf, handBuf, d.Hero, d.Weapons, mp, nil, handSize)
	if !ok {
		return TwoTurnSummary{Turn1: turn1}
	}
	return TwoTurnSummary{Turn1: turn1, Turn2: turn2}
}
