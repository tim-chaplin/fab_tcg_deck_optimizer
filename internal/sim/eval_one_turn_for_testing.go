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
	// StartOfNextTurnAuras is the live aura set at the start of the next turn: survivors
	// of this turn's chain plus any token bumps made by next-turn start-of-turn handlers.
	// Tests query specific token counts via Runechants and friends.
	StartOfNextTurnAuras []Aura
	// StartOfNextTurnItems is the live item set at the start of the next turn — items
	// carried over from this turn's chain.
	StartOfNextTurnItems []Item
	// CardsDrawn is the count of mid-chain card draws that resolved during the tested
	// turn (DrawOne calls). Surfaced for tests that want to assert on draw behaviour
	// directly without inferring it from arsenal / hand carryover.
	CardsDrawn int
	// StartOfNextTurnTriggerDamage is the damage credited by the next turn's start-of-turn
	// Aura handlers (triggers registered this turn that fired at the top of next).
	// Zero when no trigger survived. Production folds this into next turn's Value.
	StartOfNextTurnTriggerDamage int
	// StartOfNextTurnGraveyard is the auras destroyed during the next turn's start-of-turn
	// Aura pass, in destroy order.
	StartOfNextTurnGraveyard []Card
}

// Runechants returns the live Runechant token count at the start of the next turn.
func (t TurnStartState) Runechants() int {
	return tokenCountIn(t.StartOfNextTurnAuras, TokenTypeRunechant)
}

// Ponders returns the live Ponder token count at the start of the next turn.
func (t TurnStartState) Ponders() int {
	return tokenCountIn(t.StartOfNextTurnAuras, TokenTypePonder)
}

// Gold returns the live Gold token count at the start of the next turn.
func (t TurnStartState) Gold() int {
	return itemCountIn(t.StartOfNextTurnItems, TokenTypeGold)
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
	// handBuf capacity matches Evaluate's so start-of-turn Aura reveals can append
	// without realloc.
	handBuf := make([]Card, handSize, handSize+startOfTurnRevealRoom)
	tail := deckSize

	h := handBuf[:len(turn1Hand)]
	copy(h, turn1Hand)
	sortHandByID(h)
	play := Best(d.Hero, d.Weapons, h, mp, buf[head:tail], arsenalIn)
	// drawCount=0: head already points past the starting hand, so applyTurnResult only needs
	// to advance past mid-turn draws.
	nextHeld := applyTurnResult(play, buf, &head, &tail, nil)
	auraQueue := append([]Aura(nil), play.State.Auras...)
	itemQueue := append([]Item(nil), play.State.Items...)

	// Deal turn 2's hand but stop short of running Best — the caller wants the pre-Best state.
	turn2Hand, drawCount2, ok := dealNextHand(buf, handBuf, nextHeld, &head, &tail, handSize)
	if !ok {
		return TurnStartState{
			Value:                  play.Value,
			Graveyard:              append([]Card(nil), play.State.Graveyard...),
			StartOfNextTurnArsenal: play.State.Arsenal,
			StartOfNextTurnAuras:   append([]Aura(nil), play.State.Auras...),
			StartOfNextTurnItems:   append([]Item(nil), play.State.Items...),
			CardsDrawn:             play.State.CardsDrawn,
		}
	}
	// Process turn-1 Auras at the turn-2 boundary the same way Evaluate does:
	// fire start-of-turn handlers, re-arm OncePerTurn gates, drop exhausted entries.
	// Reveals into the hand are consumed here so the returned turn-2 Hand matches what
	// Best would see.
	survivors, _, trigDamage, trigRevealed, trigGraveyarded := processAurasAtStartOfTurn(auraQueue, buf[head+drawCount2:tail])
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
		StartOfNextTurnAuras:         append([]Aura(nil), survivors...),
		StartOfNextTurnItems:         append([]Item(nil), itemQueue...),
		CardsDrawn:                   play.State.CardsDrawn,
		StartOfNextTurnTriggerDamage: trigDamage,
		StartOfNextTurnGraveyard:     trigGraveyarded,
	}
}

// resolveTurn1Hand picks turn 1's starting hand and the head offset into deckCards. With
// initialHand nil the default layout takes deckCards[:handSize] as the hand and points
// head past it; with a caller-supplied hand the deck stays untouched and the supplied
// slice is used verbatim (head=0). ok=false signals the caller's inputs can't yield a
// playable opening hand: deck shorter than handSize in default mode, or a supplied hand
// that's empty or longer than handSize.
func resolveTurn1Hand(deckCards, initialHand []Card, handSize int) (hand []Card, head int, ok bool) {
	if initialHand == nil {
		if len(deckCards) < handSize {
			return nil, 0, false
		}
		return deckCards[:handSize], handSize, true
	}
	if len(initialHand) == 0 || len(initialHand) > handSize {
		return nil, 0, false
	}
	return initialHand, 0, true
}
