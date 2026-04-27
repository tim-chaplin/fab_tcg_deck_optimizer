// Package stubs provides a generic stub Hero implementation used by tests in multiple packages
// (hand, deck). It is a no-op hero — health 20, zero types, no OnCardPlayed bonus — with a
// configurable Intelligence so tests can pin hand-size-dependent behaviour without pulling in
// a real hero whose printed ability would perturb the measured value.
package stubs

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// Hero is a minimal hero.Hero. Intel is the Intelligence (hand-draw size) the stub reports;
// every other method returns a fixed zero value so the caller measures hand / deck behaviour
// in isolation from hero ability contributions.
type Hero struct {
	Intel int
}

func (Hero) ID() hero.ID                                 { return hero.Invalid }
func (Hero) Name() string                                { return "stubs.Hero" }
func (Hero) Health() int                                 { return 20 }
func (h Hero) Intelligence() int                         { return h.Intel }
func (Hero) Types() card.TypeSet                         { return 0 }
func (Hero) OnCardPlayed(card.Card, *card.TurnState) int { return 0 }
