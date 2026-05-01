// Energy Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy this: Gain {r}{r}"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var energyPotionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type EnergyPotionBlue struct{}

func (EnergyPotionBlue) ID() ids.CardID          { return ids.EnergyPotionBlue }
func (EnergyPotionBlue) Name() string            { return "Energy Potion" }
func (EnergyPotionBlue) Cost(*sim.TurnState) int { return 0 }
func (EnergyPotionBlue) Pitch() int              { return 3 }
func (EnergyPotionBlue) Attack() int             { return 0 }
func (EnergyPotionBlue) Defense() int            { return 0 }
func (EnergyPotionBlue) Types() card.TypeSet     { return energyPotionTypes }
func (EnergyPotionBlue) GoAgain() bool           { return false }

// not implemented: activated 'gain {r}{r}'
func (EnergyPotionBlue) NotImplemented()                            {}
func (EnergyPotionBlue) Play(s *sim.TurnState, self *sim.CardState) { s.LogChain(self, 0) }
