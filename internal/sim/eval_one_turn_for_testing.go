package sim

// Test-only entry point: drives one turn against a fixed deck order so tests can assert
// chain outcomes plus the start-of-next-turn state without running the full multi-shuffle
// Evaluate loop. The ForTesting suffix marks the contract; production callers use Evaluate.

// TurnStartState is the result of EvalOneTurnForTesting: the tested turn's outcome
// (Value, BestLine, Graveyard, TriggersFromLastTurn) plus a peek at the
// start-of-next-turn layout (StartOfNextTurn*) without running the next turn. For
// cross-turn behaviour (start-of-turn aura fires, banished / graveyard persistence) use
// EvalTwoTurnsForTesting.
type TurnStartState struct {
	// Value is the tested turn's TurnSummary.Value (damage dealt + damage prevented,
	// plus any start-of-turn-aura damage credited at the top of the turn).
	Value int
	// BestLine is the winning role assignment from the tested turn.
	BestLine []CardAssignment
	// Graveyard is the cards in the graveyard at end of the tested turn, in landing order.
	Graveyard []Card
	// TriggersFromLastTurn records the start-of-turn-aura fires that ran at the top of
	// this turn, each with the damage-equivalent it credited and any card it revealed
	// into the hand. Empty unless an aura was carried in via initial.Auras.
	TriggersFromLastTurn []TriggerContribution
	// StartOfNextTurnHand is the hand that next turn's dealNextHand would produce given
	// the post-turn carry — a peek, not a commit. Reflects the natural draw only;
	// next-turn start-of-turn-aura reveals are not applied.
	StartOfNextTurnHand []Card
	// StartOfNextTurnArsenal is the card in the arsenal slot at the start of the next turn.
	StartOfNextTurnArsenal Card
	// StartOfNextTurnDeck is the remaining deck after dealing next turn's hand,
	// top-to-bottom.
	StartOfNextTurnDeck []Card
	// StartOfNextTurnAuras is the live aura set at end of this turn, before any
	// start-of-turn-aura processing runs at the top of next turn.
	StartOfNextTurnAuras []Aura
	// StartOfNextTurnItems is the live item set at end of this turn — items carried
	// over from this turn's chain.
	StartOfNextTurnItems []Item
	// CardsDrawn is the count of mid-chain card draws that resolved during the tested
	// turn (DrawOne calls).
	CardsDrawn int
	// OpponentMarked is the end-of-chain Mark state on the opposing hero.
	OpponentMarked bool
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

// Silver returns the live Silver token count at the start of the next turn.
func (t TurnStartState) Silver() int {
	return itemCountIn(t.StartOfNextTurnItems, TokenTypeSilver)
}

// Copper returns the live Copper token count at the start of the next turn.
func (t TurnStartState) Copper() int {
	return itemCountIn(t.StartOfNextTurnItems, TokenTypeCopper)
}

// EvalOneTurnForTesting runs one turn against d.Cards in source order (no shuffle) and
// returns the tested turn's outcome plus a peek at the start-of-next-turn layout.
// initial seeds the start-of-turn state — Arsenal, Auras, Items, banished, graveyard,
// OpponentMarked — modelling carryover from a hypothetical previous turn; other TurnState
// fields are ignored (transient mid-chain state, hand / deck come from this method's own
// inputs). initialHand sets turn 1's starting hand; nil takes d.Cards[:handSize] as the
// hand and treats the rest as the deck, non-nil uses the slice verbatim (may be shorter
// than handSize) and treats d.Cards as the deck entirely. Test-only.
//
// initial.Auras pass through the same start-of-turn-aura processing every turn boundary
// runs, so a seeded aura with TriggerStartOfTurn fires at the top of turn 1.
func (d *Deck) EvalOneTurnForTesting(mp Matchup, initial TurnState, initialHand []Card) TurnStartState {
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

	carry := turnCarryFromInitial(initial, head, tail, handSize)
	// Stage turn1Hand into handBuf so start-of-turn-aura reveals append in-place against
	// the cushioned capacity. drawCount=0: the hand came from the caller, not from buf.
	h := handBuf[:len(turn1Hand)]
	copy(h, turn1Hand)
	art := runChainAfterDeal(&carry, buf, h, 0, d.Hero, d.Weapons, mp, nil)
	play := art.play
	carry.applyTurnAndRotate(play, buf)
	return snapshotTurnEnd(&carry, buf, play, handSize)
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

// snapshotTurnEnd builds a TurnStartState from the post-turn carry plus the chain's
// play summary. The StartOfNextTurn* fields peek at what the next dealNextHand would
// produce — they do not fire next turn's start-of-turn auras.
func snapshotTurnEnd(carry *turnCarry, buf []Card, play TurnSummary, handSize int) TurnStartState {
	state := TurnStartState{
		Value:                  play.Value,
		BestLine:               append([]CardAssignment(nil), play.BestLine...),
		Graveyard:              append([]Card(nil), play.State.Graveyard...),
		TriggersFromLastTurn:   append([]TriggerContribution(nil), play.TriggersFromLastTurn...),
		StartOfNextTurnArsenal: carry.arsenal,
		StartOfNextTurnAuras:   append([]Aura(nil), carry.auras...),
		StartOfNextTurnItems:   append([]Item(nil), carry.items...),
		CardsDrawn:             play.State.CardsDrawn,
		OpponentMarked:         carry.opponentMarked,
	}
	if hand, deck, ok := peekNextTurnLayout(buf, carry.head, carry.tail, carry.held, handSize); ok {
		state.StartOfNextTurnHand = hand
		state.StartOfNextTurnDeck = deck
	} else {
		state.StartOfNextTurnDeck = append([]Card(nil), buf[carry.head:carry.tail]...)
	}
	return state
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
