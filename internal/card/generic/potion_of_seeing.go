// Potion of Seeing — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Potion of Seeing: Look at target hero's hand."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var potionOfSeeingTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type PotionOfSeeingBlue struct{}

func (PotionOfSeeingBlue) ID() card.ID                               { return card.PotionOfSeeingBlue }
func (PotionOfSeeingBlue) Name() string                              { return "Potion of Seeing" }
func (PotionOfSeeingBlue) Cost(*card.TurnState) int                  { return 0 }
func (PotionOfSeeingBlue) Pitch() int                                { return 3 }
func (PotionOfSeeingBlue) Attack() int                               { return 0 }
func (PotionOfSeeingBlue) Defense() int                              { return 0 }
func (PotionOfSeeingBlue) Types() card.TypeSet                       { return potionOfSeeingTypes }
func (PotionOfSeeingBlue) GoAgain() bool                             { return false }
// not implemented: activated reveal opposing hero's hand
func (PotionOfSeeingBlue) NotImplemented()                           {}
func (PotionOfSeeingBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }