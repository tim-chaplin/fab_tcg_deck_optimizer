package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Item is a permanent in play with an activated ability the player chooses to play during
// their turn. Each turn the chain runner enqueues the Ability() Card as a playable option
// (1 AP, pays the printed activation cost; the Ability's Play decrements Count and removes
// the entry when Count reaches zero).
//
// Invariant: at most one Item per TokenType per TurnState; helpers bump Count on the
// existing entry rather than appending duplicates.
type Item struct {
	// Self identifies the item — a card or a token type.
	Self CardOrTokenType
	// Count is the number of copies / charges in play. The activated ability decrements
	// this each time it consumes one.
	Count int
	// Ability is the activated-ability Card the chain runner enqueues each turn. The
	// ability's Play calls back into TurnState (e.g. ConsumeItem) to decrement Count and
	// destroy this entry when Count reaches zero. Token items don't head to the graveyard
	// on destroy.
	Ability card.Card
}

// itemCountIn returns the Count of the item entry matching token type t, or zero.
func itemCountIn(items []Item, t TokenType) int {
	for i := range items {
		if items[i].Self.TokenType == t {
			return items[i].Count
		}
	}
	return 0
}
