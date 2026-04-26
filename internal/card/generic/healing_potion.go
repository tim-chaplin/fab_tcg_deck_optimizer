// Healing Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Action** - Destroy this: Gain 2{h}. **Go again**"

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var healingPotionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type HealingPotionBlue struct{}

func (HealingPotionBlue) ID() card.ID                               { return card.HealingPotionBlue }
func (HealingPotionBlue) Name() string                              { return "Healing Potion" }
func (HealingPotionBlue) Cost(*card.TurnState) int                  { return 0 }
func (HealingPotionBlue) Pitch() int                                { return 3 }
func (HealingPotionBlue) Attack() int                               { return 0 }
func (HealingPotionBlue) Defense() int                              { return 0 }
func (HealingPotionBlue) Types() card.TypeSet                       { return healingPotionTypes }
func (HealingPotionBlue) GoAgain() bool                             { return false }
// not implemented: activated 2{h} gain
func (HealingPotionBlue) NotImplemented()                           {}
func (HealingPotionBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
