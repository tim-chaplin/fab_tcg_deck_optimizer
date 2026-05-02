// Destructive Tendencies — Generic Instant. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "Choose 1 or both; - Remove all counters from target item token. - Remove all
// counters from target aura token."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var destructiveTendenciesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type DestructiveTendenciesBlue struct{}

func (DestructiveTendenciesBlue) ID() ids.CardID          { return ids.DestructiveTendenciesBlue }
func (DestructiveTendenciesBlue) Name() string            { return "Destructive Tendencies" }
func (DestructiveTendenciesBlue) Cost(*sim.TurnState) int { return 0 }
func (DestructiveTendenciesBlue) Pitch() int              { return 3 }
func (DestructiveTendenciesBlue) Attack() int             { return 0 }
func (DestructiveTendenciesBlue) Defense() int            { return 0 }
func (DestructiveTendenciesBlue) Types() card.TypeSet     { return destructiveTendenciesTypes }
func (DestructiveTendenciesBlue) GoAgain() bool           { return false }
func (DestructiveTendenciesBlue) Unplayable()             {}
func (DestructiveTendenciesBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.Log(self, 0)
}
