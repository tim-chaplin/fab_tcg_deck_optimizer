// Clarity Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Clarity Potion: **Opt 2**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var clarityPotionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type ClarityPotionBlue struct{}

func (ClarityPotionBlue) ID() ids.CardID           { return ids.ClarityPotionBlue }
func (ClarityPotionBlue) Name() string             { return "Clarity Potion" }
func (ClarityPotionBlue) Cost(*card.TurnState) int { return 0 }
func (ClarityPotionBlue) Pitch() int               { return 3 }
func (ClarityPotionBlue) Attack() int              { return 0 }
func (ClarityPotionBlue) Defense() int             { return 0 }
func (ClarityPotionBlue) Types() card.TypeSet      { return clarityPotionTypes }
func (ClarityPotionBlue) GoAgain() bool            { return false }

// not implemented: activated Opt 2
func (ClarityPotionBlue) NotImplemented()                              {}
func (ClarityPotionBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
