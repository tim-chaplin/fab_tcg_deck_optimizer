// Package hero defines the Hero interface for Flesh and Blood heroes. A deck is built around
// exactly one hero, whose class/talents gate which cards are legal and whose printed ability
// shapes the simulation.
package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Hero is a FaB hero card. Intelligence is the hand size drawn per turn; Health is starting
// life total; Types is the hero's class/talent/age set (e.g. Runeblade, Hero, Young) for O(1)
// lookup. Class is the hero's printed class (Runeblade, Brute, Thief, …) — folded into a
// Universal card's Types() so class-gated triggers (e.g. Viserai's "Runeblade card" rider)
// fire on Universal cards too. ID is the stable uint16 identifier.
type Hero interface {
	ID() ids.HeroID
	Name() string
	Health() int
	Intelligence() int
	Types() card.TypeSet
	Class() card.CardType
	// OnCardPlayed is called by the hand evaluator before each card's Play() resolves so the
	// hero's printed ability fires ahead of the card itself (matching FaB stack order). Heroes
	// that contribute damage-equivalent (e.g. a Runechant token) credit it through l (the
	// recording Logger the chain runner threads through) — hero abilities are pre-triggers
	// — and bump s.value alongside. The int return is informational and discarded by the
	// dispatcher; heroes without a triggered ability return 0.
	OnCardPlayed(played card.Card, s *TurnState, l card.Logger) int
	// Opt is the hero's heuristic for the FaB Opt N keyword. TurnState.Opt(N) pops up to N
	// cards from the top of the deck and hands them here; the handler returns a (top,
	// bottom) split. The top list is placed back on top of the deck (in returned order)
	// and the bottom list appends to the bottom (in returned order). The combined output
	// must be exactly the input multiset — adding, dropping, or substituting any card
	// panics. Both lists may be empty (skip bottoming any cards or skip keeping any on top).
	Opt(cards []card.Card) (top, bottom []card.Card)
}
