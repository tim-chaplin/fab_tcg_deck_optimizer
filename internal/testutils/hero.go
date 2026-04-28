package testutils

// Hero is a minimal heroes.Hero used by tests in multiple packages (hand, deck). It is a
// no-op hero — health 20, zero types, no OnCardPlayed bonus — with a configurable
// Intelligence so tests can pin hand-size-dependent behaviour without pulling in a real
// hero whose printed ability would perturb the measured value.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// Hero is a minimal heroes.Hero. Intel is the Intelligence (hand-draw size) the stub reports;
// every other method returns a fixed zero value so the caller measures hand / deck behaviour
// in isolation from hero ability contributions.
type Hero struct {
	Intel int
}

func (Hero) ID() ids.HeroID                              { return ids.InvalidHero }
func (Hero) Name() string                                { return "testutils.Hero" }
func (Hero) Health() int                                 { return 20 }
func (h Hero) Intelligence() int                         { return h.Intel }
func (Hero) Types() card.TypeSet                         { return 0 }
func (Hero) OnCardPlayed(card.Card, *card.TurnState) int { return 0 }
