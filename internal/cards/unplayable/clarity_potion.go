// Clarity Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Clarity Potion: **Opt 2**"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var clarityPotionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type ClarityPotionBlue struct{}

func (ClarityPotionBlue) ID() ids.CardID                             { return ids.ClarityPotionBlue }
func (ClarityPotionBlue) Name() string                               { return "Clarity Potion" }
func (ClarityPotionBlue) Cost(*sim.TurnState) int                    { return 0 }
func (ClarityPotionBlue) Pitch() int                                 { return 3 }
func (ClarityPotionBlue) Attack() int                                { return 0 }
func (ClarityPotionBlue) Defense() int                               { return 0 }
func (ClarityPotionBlue) Types() card.TypeSet                        { return clarityPotionTypes }
func (ClarityPotionBlue) GoAgain() bool                              { return false }
func (ClarityPotionBlue) Unplayable()                                {}
func (ClarityPotionBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
