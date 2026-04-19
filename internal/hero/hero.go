// Package hero defines the Hero interface for Flesh and Blood heroes. A deck is built around
// exactly one hero, whose class/talents gate which cards are legal and whose printed ability
// shapes the simulation.
package hero

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Hero is a FaB hero card. Intelligence is the hand size drawn per turn; Health is starting
// life total; Types is the hero's class/talent/age set (e.g. Runeblade, Hero, Young) for O(1)
// lookup. ID is the stable uint16 identifier used as an in-memory key (e.g. by the hand-eval
// memo).
type Hero interface {
	ID() ID
	Name() string
	Health() int
	Intelligence() int
	Types() card.TypeSet
	// OnCardPlayed is called by the hand evaluator after each card's Play() resolves and before
	// the card is appended to s.CardsPlayed. Returns any bonus damage the hero's printed ability
	// contributes (e.g. Runechant tokens). Heroes without a triggered ability return 0.
	OnCardPlayed(played card.Card, s *card.TurnState) int
}
