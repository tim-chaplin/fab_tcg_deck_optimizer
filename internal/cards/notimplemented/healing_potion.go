// Healing Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Action** - Destroy this: Gain 2{h}. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var healingPotionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type HealingPotionBlue struct{}

func (HealingPotionBlue) ID() ids.CardID          { return ids.HealingPotionBlue }
func (HealingPotionBlue) Name() string            { return "Healing Potion" }
func (HealingPotionBlue) Cost(*sim.TurnState) int { return 0 }
func (HealingPotionBlue) Pitch() int              { return 3 }
func (HealingPotionBlue) Attack() int             { return 0 }
func (HealingPotionBlue) Defense() int            { return 0 }
func (HealingPotionBlue) Types() card.TypeSet     { return healingPotionTypes }
func (HealingPotionBlue) GoAgain() bool           { return false }

// not implemented: activated 2{h} gain
func (HealingPotionBlue) NotImplemented()                            {}
func (HealingPotionBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
