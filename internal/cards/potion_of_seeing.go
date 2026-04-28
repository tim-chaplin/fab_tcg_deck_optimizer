// Potion of Seeing — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Potion of Seeing: Look at target hero's hand."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var potionOfSeeingTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type PotionOfSeeingBlue struct{}

func (PotionOfSeeingBlue) ID() ids.CardID          { return ids.PotionOfSeeingBlue }
func (PotionOfSeeingBlue) Name() string            { return "Potion of Seeing" }
func (PotionOfSeeingBlue) Cost(*sim.TurnState) int { return 0 }
func (PotionOfSeeingBlue) Pitch() int              { return 3 }
func (PotionOfSeeingBlue) Attack() int             { return 0 }
func (PotionOfSeeingBlue) Defense() int            { return 0 }
func (PotionOfSeeingBlue) Types() card.TypeSet     { return potionOfSeeingTypes }
func (PotionOfSeeingBlue) GoAgain() bool           { return false }

// not implemented: activated reveal opposing hero's hand
func (PotionOfSeeingBlue) NotImplemented()                            {}
func (PotionOfSeeingBlue) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }
