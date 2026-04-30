package testutils

// Hero is a minimal sim.Hero used by tests in multiple packages (hand, deck). It is a
// no-op hero — health 20, zero types, no OnCardPlayed bonus — with a configurable
// Intelligence so tests can pin hand-size-dependent behaviour without pulling in a real
// hero whose printed ability would perturb the measured value.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Hero is a minimal sim.Hero. Intel is the Intelligence (hand-draw size) the stub reports;
// every other method returns a fixed zero value so the caller measures hand / deck behaviour
// in isolation from hero ability contributions. OptStrategy lets tests inject a specific
// Opt heuristic (passthrough by default — every revealed card goes back on top).
type Hero struct {
	Intel       int
	OptStrategy func(cards []sim.Card) (top, bottom []sim.Card)
}

func (Hero) ID() ids.HeroID                            { return ids.InvalidHero }
func (Hero) Name() string                              { return "testutils.Hero" }
func (Hero) Health() int                               { return 20 }
func (h Hero) Intelligence() int                       { return h.Intel }
func (Hero) Types() card.TypeSet                       { return 0 }
func (Hero) OnCardPlayed(sim.Card, *sim.TurnState) int { return 0 }

// Opt dispatches to OptStrategy when set; otherwise keeps every revealed card on top of
// the deck in input order (no reshape).
func (h Hero) Opt(cards []sim.Card) (top, bottom []sim.Card) {
	if h.OptStrategy != nil {
		return h.OptStrategy(cards)
	}
	return cards, nil
}
