// Calming Breeze — Generic Instant. Cost 0. Printed pitch variants: Red 1.
//
// Text: "The next 3 times you would be dealt damage this turn, prevent 1 of that damage."
// The 3-event split is dropped — prevention is credited as a flat 3 against the next hit.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var calmingBreezeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type CalmingBreezeRed struct{}

func (CalmingBreezeRed) ID() ids.CardID          { return ids.CalmingBreezeRed }
func (CalmingBreezeRed) Name() string            { return "Calming Breeze" }
func (CalmingBreezeRed) Cost(*sim.TurnState) int { return 0 }
func (CalmingBreezeRed) Pitch() int              { return 1 }
func (CalmingBreezeRed) Attack() int             { return 0 }
func (CalmingBreezeRed) Defense() int            { return 3 }
func (CalmingBreezeRed) Types() card.TypeSet     { return calmingBreezeTypes }
func (CalmingBreezeRed) GoAgain() bool           { return false }
func (CalmingBreezeRed) DefensiveInstant()       {}
func (CalmingBreezeRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}
