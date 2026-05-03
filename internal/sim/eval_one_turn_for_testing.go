package sim

// Test-only entry point: drives one turn against a fixed deck order so tests can assert
// chain outcomes plus the start-of-next-turn state without running the full multi-shuffle
// Evaluate loop. The ForTesting suffix marks the contract; production callers use Evaluate.

// TurnStartState is the result of EvalOneTurnForTesting: the tested turn's outcome (Value,
// BestLine, Graveyard) plus the start-of-next-turn state (StartOfNextTurn*) so cross-turn
// effects can be asserted without simulating the next turn.
type TurnStartState struct {
	// Value is the tested turn's TurnSummary.Value (damage dealt + damage prevented).
	Value int
	// BestLine is the winning role assignment from the tested turn.
	BestLine []CardAssignment
	// Graveyard is the cards in the graveyard at end of the tested turn, in landing order.
	// Distinguishes "in graveyard" from "absent from next-turn surfaces".
	Graveyard []Card
	// StartOfNextTurnHand is the hand dealt for the turn after the tested turn.
	StartOfNextTurnHand []Card
	// StartOfNextTurnArsenal is the card in the arsenal slot at the start of the next turn.
	StartOfNextTurnArsenal Card
	// StartOfNextTurnDeck is the remaining deck at the start of the next turn, top-to-bottom.
	StartOfNextTurnDeck []Card
	// StartOfNextTurnRunechants is the live Runechant count at the start of the next turn:
	// chain leftover plus tokens created by next-turn start-of-turn AuraTrigger handlers.
	StartOfNextTurnRunechants int
	// StartOfNextTurnTriggerDamage is the damage credited by the next turn's start-of-turn
	// AuraTrigger handlers (triggers registered this turn that fired at the top of next).
	// Zero when no trigger survived. Production folds this into next turn's Value.
	StartOfNextTurnTriggerDamage int
	// StartOfNextTurnGraveyard is the auras destroyed during the next turn's start-of-turn
	// AuraTrigger pass, in destroy order.
	StartOfNextTurnGraveyard []Card
}

// EvalOneTurnForTesting runs one turn against d.Cards in source order (no shuffle) and returns
// the tested turn's outcome plus the start-of-next-turn state. arsenalIn seeds turn 1's arsenal
// slot (nil for empty). initialHand sets turn 1's starting hand; nil takes d.Cards[:handSize] as
// the hand and treats the rest as the deck, non-nil uses the slice directly (may be shorter than
// handSize) and treats d.Cards as the deck entirely. Test-only — production callers use Evaluate.
func (d *Deck) EvalOneTurnForTesting(mp Matchup, arsenalIn Card, initialHand []Card) TurnStartState {
	CurrentHero = d.Hero
	handSize := d.Hero.Intelligence()
	if handSize <= 0 {
		return TurnStartState{}
	}

	turn1Hand, head, ok := resolveTurn1Hand(d.Cards, initialHand, handSize)
	if !ok {
		return TurnStartState{}
	}

	deckSize := len(d.Cards)
	// Oversized buf: 2×deckSize matches Evaluate's layout. Add a handSize cushion so small
	// decks still have room for mid-turn pitches (hand + drawn) without overflowing tail.
	buf := make([]Card, deckSize*2+handSize*2)
	copy(buf, d.Cards)
	// handBuf capacity matches Evaluate's so start-of-turn AuraTrigger reveals can append
	// without realloc.
	handBuf := make([]Card, handSize, handSize+startOfTurnRevealRoom)
	tail := deckSize

	h := handBuf[:len(turn1Hand)]
	copy(h, turn1Hand)
	sortHandByID(h)
	play := Best(d.Hero, d.Weapons, h, mp, buf[head:tail], 0, arsenalIn)
	// drawCount=0: head already points past the starting hand, so applyTurnResult only needs
	// to advance past mid-turn draws.
	nextHeld := applyTurnResult(play, buf, &head, &tail, nil)
	triggerQueue := append([]AuraTrigger(nil), play.State.AuraTriggers...)

	// Deal turn 2's hand but stop short of running Best — the caller wants the pre-Best state.
	turn2Hand, drawCount2, ok := dealNextHand(buf, handBuf, nextHeld, &head, &tail, handSize)
	if !ok {
		return TurnStartState{
			Value:                     play.Value,
			Graveyard:                 append([]Card(nil), play.State.Graveyard...),
			StartOfNextTurnArsenal:    play.State.Arsenal,
			StartOfNextTurnRunechants: play.State.Runechants,
		}
	}
	// Process turn-1 AuraTriggers at the turn-2 boundary the same way Evaluate does:
	// fire start-of-turn handlers, re-arm OncePerTurn gates, drop exhausted entries.
	// Reveals into the hand are consumed here so the returned turn-2 Hand matches what
	// Best would see.
	_, _, trigDamage, trigRunes, trigRevealed, trigGraveyarded := processTriggersAtStartOfTurn(triggerQueue, buf[head+drawCount2:tail])
	for range trigRevealed {
		turn2Hand = append(turn2Hand, buf[head+drawCount2])
		drawCount2++
	}
	handCopy := append([]Card(nil), turn2Hand...)
	deckLeft := append([]Card(nil), buf[head+drawCount2:tail]...)
	lineCopy := append([]CardAssignment(nil), play.BestLine...)

	return TurnStartState{
		Value:                        play.Value,
		BestLine:                     lineCopy,
		Graveyard:                    append([]Card(nil), play.State.Graveyard...),
		StartOfNextTurnHand:          handCopy,
		StartOfNextTurnArsenal:       play.State.Arsenal,
		StartOfNextTurnDeck:          deckLeft,
		StartOfNextTurnRunechants:    play.State.Runechants + trigRunes,
		StartOfNextTurnTriggerDamage: trigDamage,
		StartOfNextTurnGraveyard:     trigGraveyarded,
	}
}
