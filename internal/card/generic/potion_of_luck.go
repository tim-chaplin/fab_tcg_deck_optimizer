// Potion of Luck — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Potion of Luck: Shuffle your hand and arsenal into your deck then
// draw that many cards."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var potionOfLuckTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type PotionOfLuckBlue struct{}

func (PotionOfLuckBlue) ID() card.ID                               { return card.PotionOfLuckBlue }
func (PotionOfLuckBlue) Name() string                              { return "Potion of Luck" }
func (PotionOfLuckBlue) Cost(*card.TurnState) int                  { return 0 }
func (PotionOfLuckBlue) Pitch() int                                { return 3 }
func (PotionOfLuckBlue) Attack() int                               { return 0 }
func (PotionOfLuckBlue) Defense() int                              { return 0 }
func (PotionOfLuckBlue) Types() card.TypeSet                       { return potionOfLuckTypes }
func (PotionOfLuckBlue) GoAgain() bool                             { return false }
// not implemented: activated 'shuffle hand+arsenal into deck, draw that many'
func (PotionOfLuckBlue) NotImplemented()                           {}
func (PotionOfLuckBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }