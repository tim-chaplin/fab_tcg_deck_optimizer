package sim

// Test-only entry point: drives one turn against a fixed deck order so tests can assert
// chain outcomes plus the start-of-next-turn state without running the full multi-shuffle
// Evaluate loop. The ForTesting suffix marks the contract; production callers use Evaluate.

// EvalOneTurnForTesting runs one turn against d.Cards in source order (no shuffle) and
// returns the chain's TurnSummary. initial seeds the start-of-turn state — Arsenal,
// Auras, Items, banished, graveyard, OpponentMarked — modelling carryover from a
// hypothetical previous turn; other TurnState fields are ignored (transient mid-chain
// state, hand / deck come from this method's own inputs). initialHand sets turn 1's
// starting hand; nil takes d.Cards[:handSize] as the hand and treats the rest as the
// deck, non-nil uses the slice verbatim (may be shorter than handSize) and treats
// d.Cards as the deck entirely. Test-only.
//
// initial.Auras pass through the same start-of-turn-aura processing every turn boundary
// runs, so a seeded aura with TriggerStartOfTurn fires at the top of turn 1.
func (d *Deck) EvalOneTurnForTesting(mp Matchup, initial TurnState, initialHand []Card) TurnSummary {
	CurrentHero = d.Hero
	handSize := d.Hero.Intelligence()
	if handSize <= 0 {
		return TurnSummary{}
	}

	turn1Hand, head, ok := resolveTurn1Hand(d.Cards, initialHand, handSize)
	if !ok {
		return TurnSummary{}
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

	carry := turnCarryFromInitial(initial, head, tail, handSize)
	// Stage turn1Hand into handBuf so start-of-turn-aura reveals append in-place against
	// the cushioned capacity. drawCount=0: the hand came from the caller, not from buf.
	h := handBuf[:len(turn1Hand)]
	copy(h, turn1Hand)
	art := runChainAfterDeal(&carry, buf, h, 0, d.Hero, d.Weapons, mp, nil)
	return art.play
}

// turnCarryFromInitial builds a turnCarry seeded with caller-supplied prior-turn state.
// Held starts empty; callers that supply turn 1's hand pass it directly to
// runChainAfterDeal. Allocates fresh backing for every slice so the caller's inputs
// aren't aliased.
func turnCarryFromInitial(initial TurnState, head, tail, handSize int) turnCarry {
	return turnCarry{
		head:           head,
		tail:           tail,
		arsenal:        initial.Arsenal,
		opponentMarked: initial.OpponentMarked,
		held:           make([]Card, 0, handSize),
		nextHeld:       make([]Card, 0, handSize),
		auras:          append([]Aura(nil), initial.Auras...),
		nextAuras:      make([]Aura, 0, len(initial.Auras)+4),
		items:          append([]Item(nil), initial.Items...),
		nextItems:      make([]Item, 0, len(initial.Items)+4),
		banish:         append([]Card(nil), initial.banished...),
		nextBanish:     make([]Card, 0, len(initial.banished)+4),
		graveyard:      append([]Card(nil), initial.graveyard...),
		nextGraveyard:  make([]Card, 0, len(initial.graveyard)+4),
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
