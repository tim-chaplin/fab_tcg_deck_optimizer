// Potion of Déjà Vu — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Potion of Déjà Vu: Put all cards from your pitch zone on top of your
// deck in any order."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var potionOfDejaVuTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type PotionOfDejaVuBlue struct{}

func (PotionOfDejaVuBlue) ID() card.ID                               { return card.PotionOfDejaVuBlue }
func (PotionOfDejaVuBlue) Name() string                              { return "Potion of Déjà Vu" }
func (PotionOfDejaVuBlue) Cost(*card.TurnState) int                  { return 0 }
func (PotionOfDejaVuBlue) Pitch() int                                { return 3 }
func (PotionOfDejaVuBlue) Attack() int                               { return 0 }
func (PotionOfDejaVuBlue) Defense() int                              { return 0 }
func (PotionOfDejaVuBlue) Types() card.TypeSet                       { return potionOfDejaVuTypes }
func (PotionOfDejaVuBlue) GoAgain() bool                             { return false }
// not implemented: activated 'put pitch zone on top of deck in any order'
func (PotionOfDejaVuBlue) NotImplemented()                           {}
func (PotionOfDejaVuBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }