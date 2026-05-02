// Potion of Strength — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Action** - Destroy this: Your next attack this turn gains +2{p}. **Go again**"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var potionOfStrengthTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type PotionOfStrengthBlue struct{}

func (PotionOfStrengthBlue) ID() ids.CardID                             { return ids.PotionOfStrengthBlue }
func (PotionOfStrengthBlue) Name() string                               { return "Potion of Strength" }
func (PotionOfStrengthBlue) Cost(*sim.TurnState) int                    { return 0 }
func (PotionOfStrengthBlue) Pitch() int                                 { return 3 }
func (PotionOfStrengthBlue) Attack() int                                { return 0 }
func (PotionOfStrengthBlue) Defense() int                               { return 0 }
func (PotionOfStrengthBlue) Types() card.TypeSet                        { return potionOfStrengthTypes }
func (PotionOfStrengthBlue) GoAgain() bool                              { return false }
func (PotionOfStrengthBlue) Unplayable()                                {}
func (PotionOfStrengthBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
