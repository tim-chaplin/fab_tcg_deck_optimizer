// Blade Flash — Generic Attack Reaction. Cost 1. Printed pitch variants: Blue 3. Defense 2.
//
// Text: "Target sword attack gains **go again**."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var bladeFlashTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type BladeFlashBlue struct{}

func (BladeFlashBlue) ID() card.ID                               { return card.BladeFlashBlue }
func (BladeFlashBlue) Name() string                              { return "Blade Flash" }
func (BladeFlashBlue) Cost(*card.TurnState) int                  { return 1 }
func (BladeFlashBlue) Pitch() int                                { return 3 }
func (BladeFlashBlue) Attack() int                               { return 0 }
func (BladeFlashBlue) Defense() int                              { return 2 }
func (BladeFlashBlue) Types() card.TypeSet                       { return bladeFlashTypes }
func (BladeFlashBlue) GoAgain() bool                             { return false }
// not implemented: AR 'target sword attack gains go again'
func (BladeFlashBlue) NotImplemented()                           {}
func (BladeFlashBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }