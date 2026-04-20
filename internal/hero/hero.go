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

// PlayTypeFilter is an optional narrow contract: heroes whose OnCardPlayed trigger gates on the
// played card's TypeSet can implement it to let the solver short-circuit the interface call on
// clearly-irrelevant cards. The solver has each card's cached TypeSet available on attackerMeta,
// so routing the pre-check through this method avoids a per-play played.Types() interface
// dispatch inside OnCardPlayed for cards that would return 0 anyway. Heroes whose ability
// doesn't key on the played card's type alone (e.g. triggers reading deck contents or other
// state) should leave this unimplemented — the solver will simply call OnCardPlayed
// unconditionally.
type PlayTypeFilter interface {
	// CardTypeCanTrigger reports whether a card with the given TypeSet might cause this hero's
	// OnCardPlayed to return a non-zero bonus. False lets the solver skip OnCardPlayed entirely.
	// Must be conservative: returning true when the hero ends up returning 0 is a missed
	// optimization but correct; returning false when the hero would have returned >0 is a bug.
	CardTypeCanTrigger(card.TypeSet) bool
}
