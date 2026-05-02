// Potion of Luck — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Potion of Luck: Shuffle your hand and arsenal into your deck then
// draw that many cards."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var potionOfLuckTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type PotionOfLuckBlue struct{}

func (PotionOfLuckBlue) ID() ids.CardID                             { return ids.PotionOfLuckBlue }
func (PotionOfLuckBlue) Name() string                               { return "Potion of Luck" }
func (PotionOfLuckBlue) Cost(*sim.TurnState) int                    { return 0 }
func (PotionOfLuckBlue) Pitch() int                                 { return 3 }
func (PotionOfLuckBlue) Attack() int                                { return 0 }
func (PotionOfLuckBlue) Defense() int                               { return 0 }
func (PotionOfLuckBlue) Types() card.TypeSet                        { return potionOfLuckTypes }
func (PotionOfLuckBlue) GoAgain() bool                              { return false }
func (PotionOfLuckBlue) Unplayable()                                {}
func (PotionOfLuckBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
