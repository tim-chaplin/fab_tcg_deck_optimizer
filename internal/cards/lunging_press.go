// Lunging Press — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3. Defense 2.
//
// Text: "Target attack action card gains +1{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var lungingPressTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type LungingPressBlue struct{}

func (LungingPressBlue) ID() ids.CardID           { return ids.LungingPressBlue }
func (LungingPressBlue) Name() string             { return "Lunging Press" }
func (LungingPressBlue) Cost(*card.TurnState) int { return 0 }
func (LungingPressBlue) Pitch() int               { return 3 }
func (LungingPressBlue) Attack() int              { return 0 }
func (LungingPressBlue) Defense() int             { return 2 }
func (LungingPressBlue) Types() card.TypeSet      { return lungingPressTypes }
func (LungingPressBlue) GoAgain() bool            { return false }

// not implemented: AR +1{p} buff to a target attack action card
func (LungingPressBlue) NotImplemented()                              {}
func (LungingPressBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
