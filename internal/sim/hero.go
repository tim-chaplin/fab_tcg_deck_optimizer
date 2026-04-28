// Package hero defines the Hero interface for Flesh and Blood heroes. A deck is built around
// exactly one hero, whose class/talents gate which cards are legal and whose printed ability
// shapes the simulation.
package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// Hero is a FaB hero card. Intelligence is the hand size drawn per turn; Health is starting
// life total; Types is the hero's class/talent/age set (e.g. Runeblade, Hero, Young) for O(1)
// lookup. ID is the stable uint16 identifier.
type Hero interface {
	ID() ids.HeroID
	Name() string
	Health() int
	Intelligence() int
	Types() card.TypeSet
	// OnCardPlayed is called by the hand evaluator before each card's Play() resolves so the
	// hero's printed ability fires ahead of the card itself (matching FaB stack order). Heroes
	// that contribute damage-equivalent (e.g. a Runechant token) credit it through
	// s.AddPreTriggerLogEntry — hero abilities are pre-triggers — which both writes the
	// trigger's log line and bumps s.Value. The int return is informational and discarded by
	// the dispatcher — it's the value AddPreTriggerLogEntry already credited, surfaced so
	// callers can fold the call into a single return statement. Heroes without a triggered
	// ability return 0.
	OnCardPlayed(played Card, s *TurnState) int
}
