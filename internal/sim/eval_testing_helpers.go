package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// snapshotState clones src's persistent carryover and overlays deck (copied), hand, value,
// and cardsDrawn. cardsDrawn is passed explicitly because ResetEphemeralState wipes it on
// src before the snapshot is taken.
func snapshotState(src *gameengine.GameState, d *deck.Deck, hand []card.Card, value, cardsDrawn int) *gameengine.GameState {
	out := src.CopyPersistentState()
	if d != nil {
		out.SetDeck(d.Copy())
	}
	out.SetHand(hand)
	out.SetValue(value)
	out.SetCardsDrawn(cardsDrawn)
	return out
}

// cardsToDeckCards widens []card.Card to []deck.Card; returns nil for empty input.
func cardsToDeckCards(in []card.Card) []deck.Card {
	if len(in) == 0 {
		return nil
	}
	out := make([]deck.Card, len(in))
	for i, c := range in {
		out[i] = c
	}
	return out
}
