// Energy Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy this: Gain {r}{r}"

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var energyPotionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type EnergyPotionBlue struct{}

func (EnergyPotionBlue) ID() card.ID                               { return card.EnergyPotionBlue }
func (EnergyPotionBlue) Name() string                              { return "Energy Potion" }
func (EnergyPotionBlue) Cost(*card.TurnState) int                  { return 0 }
func (EnergyPotionBlue) Pitch() int                                { return 3 }
func (EnergyPotionBlue) Attack() int                               { return 0 }
func (EnergyPotionBlue) Defense() int                              { return 0 }
func (EnergyPotionBlue) Types() card.TypeSet                       { return energyPotionTypes }
func (EnergyPotionBlue) GoAgain() bool                             { return false }
// not implemented: activated 'gain {r}{r}'
func (EnergyPotionBlue) NotImplemented()                           {}
func (EnergyPotionBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }