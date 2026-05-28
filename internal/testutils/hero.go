package testutils

// Hero is a minimal gameengine.Hero used by tests in multiple packages (hand, deck). It is
// a no-op hero — health 20, zero types, no triggered ability — with a configurable
// Intelligence so tests can pin hand-size-dependent behaviour without pulling in a real
// hero whose printed ability would perturb the measured value.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Hero is a minimal gameengine.Hero. Intel is the Intelligence (hand-draw size) the stub
// reports; every other method returns a fixed zero value so the caller measures hand /
// deck behaviour in isolation from hero ability contributions. OptStrategy lets tests
// inject a specific Opt heuristic (passthrough by default — every revealed card goes back
// on top).
type Hero struct {
	Intel       int
	OptStrategy func(cards []card.Card) (top, bottom []card.Card)
	TypeSet     card.TypeSet
}

func (Hero) ID() ids.HeroID        { return ids.InvalidHero }
func (Hero) Name() string          { return "testutils.Hero" }
func (Hero) DisplayName() string   { return "testutils.Hero" }
func (h Hero) Intelligence() int   { return h.Intel }
func (h Hero) Types() card.TypeSet { return h.TypeSet }
func (Hero) Class() card.CardType  { return 0 }

// No triggered ability — TriggerType == 0 makes FireTriggers skip the hero.
func (Hero) TriggerType() triggertype.Type                                        { return 0 }
func (Hero) OncePerTurn() bool                                                    { return false }
func (Hero) FiredThisTurn() bool                                                  { return false }
func (Hero) SetFiredThisTurn(bool)                                                {}
func (Hero) Matches(card.TypeSet) bool                                            { return false }
func (Hero) Fire(card.GameEngine, card.Logger, *card.CardState, triggertype.Type) {}

// Opt dispatches to OptStrategy when set; otherwise keeps every revealed card on top of
// the deck in input order (no reshape).
func (h Hero) Opt(cards []card.Card) (top, bottom []card.Card) {
	if h.OptStrategy != nil {
		return h.OptStrategy(cards)
	}
	return cards, nil
}
