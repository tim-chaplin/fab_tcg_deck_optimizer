// Blade Flash — Generic Attack Reaction. Cost 1. Printed pitch variants: Blue 3. Defense 2.
//
// Text: "Target sword attack gains **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var bladeFlashTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type BladeFlashBlue struct{}

func (BladeFlashBlue) ID() ids.CardID           { return ids.BladeFlashBlue }
func (BladeFlashBlue) Name() string             { return "Blade Flash" }
func (BladeFlashBlue) Cost(*card.TurnState) int { return 1 }
func (BladeFlashBlue) Pitch() int               { return 3 }
func (BladeFlashBlue) Attack() int              { return 0 }
func (BladeFlashBlue) Defense() int             { return 2 }
func (BladeFlashBlue) Types() card.TypeSet      { return bladeFlashTypes }
func (BladeFlashBlue) GoAgain() bool            { return false }

// not implemented: AR 'target sword attack gains go again'
func (BladeFlashBlue) NotImplemented()                              {}
func (BladeFlashBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
