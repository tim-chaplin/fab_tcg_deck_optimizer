// Potion of Strength — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Action** - Destroy this: Your next attack this turn gains +2{p}. **Go again**"

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var potionOfStrengthTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type PotionOfStrengthBlue struct{}

func (PotionOfStrengthBlue) ID() card.ID              { return card.PotionOfStrengthBlue }
func (PotionOfStrengthBlue) Name() string             { return "Potion of Strength" }
func (PotionOfStrengthBlue) Cost(*card.TurnState) int { return 0 }
func (PotionOfStrengthBlue) Pitch() int               { return 3 }
func (PotionOfStrengthBlue) Attack() int              { return 0 }
func (PotionOfStrengthBlue) Defense() int             { return 0 }
func (PotionOfStrengthBlue) Types() card.TypeSet      { return potionOfStrengthTypes }
func (PotionOfStrengthBlue) GoAgain() bool            { return false }

// not implemented: activated +2{p} on next attack
func (PotionOfStrengthBlue) NotImplemented()                              {}
func (PotionOfStrengthBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
