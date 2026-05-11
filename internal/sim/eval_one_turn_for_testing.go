package sim

import "github.com/tim-chaplin/fab-deck-optimizer/v2/deck"

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
	Graveyard []deck.Card
	// StartOfNextTurnHand is the hand dealt for the turn after the tested turn.
	StartOfNextTurnHand []deck.Card
	// StartOfNextTurnArsenal is the card in the arsenal slot at the start of the next turn.
	StartOfNextTurnArsenal deck.Card
	// StartOfNextTurnDeck is the remaining deck at the start of the next turn — a
	// fresh *deck.Deck the caller can drive through Draw / PeekTop / Size.
	StartOfNextTurnDeck *deck.Deck
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
	// OpponentMarked is the end-of-chain Mark state on the opposing hero — true when
	// the winning chain landed (and didn't subsequently strip) a Mark.
	OpponentMarked bool
	// StartOfNextTurnTriggerDamage is the damage credited by the next turn's start-of-turn
	// Aura handlers (triggers registered this turn that fired at the top of next).
	// Zero when no trigger survived. Production folds this into next turn's Value.
	StartOfNextTurnTriggerDamage int
	// StartOfNextTurnGraveyard is the auras destroyed during the next turn's start-of-turn
	// Aura pass, in destroy order.
	StartOfNextTurnGraveyard []deck.Card
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

// EvalOneTurnForTesting runs one turn against the deck in source order (no shuffle) and returns
// the tested turn's outcome plus the start-of-next-turn state. initial seeds the start-of-turn
// state — Arsenal, Auras, Items — modelling carryover from a hypothetical previous turn; the
// other TurnState fields are ignored (transient mid-chain state, hand / deck / graveyard which
// are seeded from this function's own inputs). External tests build initial via
// NewTurnStateFromSpec(TurnStateSpec{...}). initialHand sets turn 1's starting hand; nil draws
// the hand off the top of the deck, non-nil uses the slice directly (may be shorter than
// handSize) and leaves the deck untouched. Test-only — production callers use Evaluate.
//
// Free function (not a method) because deck.Deck lives in another package; Go disallows
// methods on imported types.
func EvalOneTurnForTesting(master *deck.Deck, mp Matchup, initial TurnState, initialHand []deck.Card) TurnStartState {
	d := master.Copy()
	hero := d.Hero.(Hero)
	CurrentHero = hero
	handSize := hero.Intelligence()
	if handSize <= 0 {
		return TurnStartState{}
	}

	// Build turn 1's hand. Caller-supplied initialHand is used verbatim and the deck stays
	// untouched; otherwise draw handSize cards off the top.
	if initialHand == nil && d.Size() < handSize {
		return TurnStartState{}
	}
	if initialHand != nil && (len(initialHand) == 0 || len(initialHand) > handSize) {
		return TurnStartState{}
	}
	handBuf := make([]Card, handSize, handSize+startOfTurnRevealRoom)
	var h []Card
	if initialHand != nil {
		h = handBuf[:len(initialHand)]
		for i, c := range initialHand {
			h[i] = c.(Card)
		}
	} else {
		h = handBuf[:handSize]
		for i, c := range d.Draw(handSize) {
			h[i] = c.(Card)
		}
	}
	sortHandByID(h)
	weapons := make([]Weapon, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = w.(Weapon)
	}
	play := best(hero, weapons, h, mp, d, initial)
	// Adopt the chain's post-mutation deck and recycle pitched cards onto the bottom.
	d = play.State.Deck
	pitched := pitchedFromBestLine(play.BestLine)
	recycled := make([]deck.Card, len(pitched))
	for i, c := range pitched {
		recycled[i] = c
	}
	d.PutBottom(recycled)

	auraQueue := append([]Aura(nil), play.State.Auras...)
	itemQueue := append([]Item(nil), play.State.Items...)

	// Deal turn 2's hand off the top, with the chain's leftover hand as the held prefix.
	// Stop short of running Best — the caller wants the pre-Best state.
	held := append([]Card(nil), play.State.Hand...)
	graveyardOut := make([]deck.Card, len(play.State.Graveyard))
	for i, c := range play.State.Graveyard {
		graveyardOut[i] = c
	}
	if len(held) >= handSize || d.Size() < handSize-len(held) {
		return TurnStartState{
			Value:                  play.Value,
			BestLine:               append([]CardAssignment(nil), play.BestLine...),
			Graveyard:              graveyardOut,
			StartOfNextTurnArsenal: play.State.Arsenal,
			StartOfNextTurnDeck:    d.Copy(),
			StartOfNextTurnAuras:   append([]Aura(nil), play.State.Auras...),
			StartOfNextTurnItems:   append([]Item(nil), play.State.Items...),
			CardsDrawn:             play.State.CardsDrawn,
			OpponentMarked:         play.State.OpponentMarked,
		}
	}
	turn2Hand := append([]Card(nil), held...)
	for _, c := range d.Draw(handSize - len(held)) {
		turn2Hand = append(turn2Hand, c.(Card))
	}
	// Process turn-1 Auras at the turn-2 boundary the same way Evaluate does:
	// fire start-of-turn handlers, re-arm OncePerTurn gates, drop exhausted entries.
	// Reveals into the hand are consumed here so the returned turn-2 Hand matches what
	// Best would see.
	survivors, _, trigDamage, trigRevealed, trigGraveyarded := processAurasAtStartOfTurn(auraQueue, d)
	for _, c := range trigRevealed {
		turn2Hand = append(turn2Hand, c)
	}
	handOut := make([]deck.Card, len(turn2Hand))
	for i, c := range turn2Hand {
		handOut[i] = c
	}
	graveyardedOut := make([]deck.Card, len(trigGraveyarded))
	for i, c := range trigGraveyarded {
		graveyardedOut[i] = c
	}

	return TurnStartState{
		Value:                        play.Value,
		BestLine:                     append([]CardAssignment(nil), play.BestLine...),
		Graveyard:                    graveyardOut,
		StartOfNextTurnHand:          handOut,
		StartOfNextTurnArsenal:       play.State.Arsenal,
		StartOfNextTurnDeck:          d.Copy(),
		StartOfNextTurnAuras:         append([]Aura(nil), survivors...),
		StartOfNextTurnItems:         append([]Item(nil), itemQueue...),
		CardsDrawn:                   play.State.CardsDrawn,
		OpponentMarked:               play.State.OpponentMarked,
		StartOfNextTurnTriggerDamage: trigDamage,
		StartOfNextTurnGraveyard:     graveyardedOut,
	}
}
