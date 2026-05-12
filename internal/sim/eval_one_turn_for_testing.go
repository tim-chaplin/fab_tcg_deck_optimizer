package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// Test-only entry point: drives one turn against a fixed deck order so tests can assert
// chain outcomes plus the start-of-next-turn state without running the full multi-shuffle
// Evaluate loop. The ForTesting suffix marks the contract; production callers use Evaluate.

// TurnStartState is the result of EvalOneTurnForTesting: the tested turn's outcome (Value,
// BestLine, Graveyard) plus the start-of-next-turn state (StartOfNextTurn*) so cross-turn
// effects can be asserted without simulating the next turn.
type TurnStartState struct {
	Value                        int
	BestLine                     []CardAssignment
	Graveyard                    []deck.Card
	StartOfNextTurnHand          []deck.Card
	StartOfNextTurnArsenal       deck.Card
	StartOfNextTurnDeck          *deck.Deck
	StartOfNextTurnAuras         []*Aura
	StartOfNextTurnItems         []*Item
	CardsDrawn                   int
	OpponentMarked               bool
	StartOfNextTurnTriggerDamage int
	StartOfNextTurnGraveyard     []deck.Card
}

// RunechantCount returns the live Runechant token count at the start of the next turn.
func (t TurnStartState) RunechantCount() int {
	return auraCountByName(t.StartOfNextTurnAuras, "Runechant")
}

// PonderCount returns the live Ponder token count at the start of the next turn.
func (t TurnStartState) PonderCount() int {
	return auraCountByName(t.StartOfNextTurnAuras, "Ponder")
}

// GoldCount returns the live Gold token count at the start of the next turn.
func (t TurnStartState) GoldCount() int {
	return itemCountByName(t.StartOfNextTurnItems, "Gold")
}

// SilverCount returns the live Silver token count at the start of the next turn.
func (t TurnStartState) SilverCount() int {
	return itemCountByName(t.StartOfNextTurnItems, "Silver")
}

// CopperCount returns the live Copper token count at the start of the next turn.
func (t TurnStartState) CopperCount() int {
	return itemCountByName(t.StartOfNextTurnItems, "Copper")
}

// EvalOneTurnForTesting runs one turn against the deck in source order (no shuffle) and
// returns the tested turn's outcome plus the start-of-next-turn state. initial seeds the
// turn's prior carryover (Auras / Items / Arsenal / Banished / Graveyard /
// OpponentMarked); leave its Hero nil to inherit the deck's hero.
func EvalOneTurnForTesting(master *deck.Deck, mp Matchup, initial Prior, initialHand []deck.Card) TurnStartState {
	d := master.Copy()
	hero := d.Hero.(Hero)
	if initial.Hero == nil {
		initial.Hero = hero
	} else {
		hero = initial.Hero
	}
	handSize := hero.Intelligence()
	if handSize <= 0 {
		return TurnStartState{}
	}
	if initialHand == nil && d.Size() < handSize {
		return TurnStartState{}
	}
	if initialHand != nil && (len(initialHand) == 0 || len(initialHand) > handSize) {
		return TurnStartState{}
	}
	handBuf := make([]card.Card, handSize, handSize+startOfTurnRevealRoom)
	var h []card.Card
	if initialHand != nil {
		h = handBuf[:len(initialHand)]
		for i, c := range initialHand {
			h[i] = c.(card.Card)
		}
	} else {
		h = handBuf[:handSize]
		for i, c := range d.Draw(handSize) {
			h[i] = c.(card.Card)
		}
	}
	sortHandByID(h)
	weapons := make([]Weapon, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = w.(Weapon)
	}
	play := best(weapons, h, mp, d, initial)
	d = play.State.DeckRaw()
	pitched := pitchedFromBestLine(play.BestLine)
	recycled := make([]deck.Card, len(pitched))
	for i, c := range pitched {
		recycled[i] = c
	}
	d.PutBottom(recycled)

	auraQueue := make([]*Aura, 0, len(play.State.Auras()))
	for _, a := range play.State.Auras() {
		auraQueue = append(auraQueue, a.(*Aura))
	}
	itemQueue := make([]*Item, 0, len(play.State.Items()))
	for _, it := range play.State.Items() {
		itemQueue = append(itemQueue, it.(*Item))
	}

	held := append([]card.Card(nil), play.State.HandRaw()...)
	graveyardOut := make([]deck.Card, len(play.State.GraveyardRaw()))
	for i, c := range play.State.GraveyardRaw() {
		graveyardOut[i] = c
	}
	if len(held) >= handSize || d.Size() < handSize-len(held) {
		return TurnStartState{
			Value:                  play.Value,
			BestLine:               append([]CardAssignment(nil), play.BestLine...),
			Graveyard:              graveyardOut,
			StartOfNextTurnArsenal: play.State.Arsenal(),
			StartOfNextTurnDeck:    d.Copy(),
			StartOfNextTurnAuras:   append([]*Aura(nil), auraQueue...),
			StartOfNextTurnItems:   append([]*Item(nil), itemQueue...),
			CardsDrawn:             play.State.CardsDrawn(),
			OpponentMarked:         play.State.OpponentMarked(),
		}
	}
	turn2Hand := append([]card.Card(nil), held...)
	for _, c := range d.Draw(handSize - len(held)) {
		turn2Hand = append(turn2Hand, c.(card.Card))
	}
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
		StartOfNextTurnArsenal:       play.State.Arsenal(),
		StartOfNextTurnDeck:          d.Copy(),
		StartOfNextTurnAuras:         append([]*Aura(nil), survivors...),
		StartOfNextTurnItems:         append([]*Item(nil), itemQueue...),
		CardsDrawn:                   play.State.CardsDrawn(),
		OpponentMarked:               play.State.OpponentMarked(),
		StartOfNextTurnTriggerDamage: trigDamage,
		StartOfNextTurnGraveyard:     graveyardedOut,
	}
}
