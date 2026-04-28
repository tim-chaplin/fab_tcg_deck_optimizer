// Cracked Bauble — Generic Resource. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "*(A player may add any number of Cracked Baubles to their card-pool in sealed deck or
// booster draft formats.)*"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var crackedBaubleTypes = card.NewTypeSet(card.TypeGeneric)

type CrackedBaubleYellow struct{}

func (CrackedBaubleYellow) ID() ids.CardID           { return ids.CrackedBaubleYellow }
func (CrackedBaubleYellow) Name() string             { return "Cracked Bauble" }
func (CrackedBaubleYellow) Cost(*card.TurnState) int { return 0 }
func (CrackedBaubleYellow) Pitch() int               { return 2 }
func (CrackedBaubleYellow) Attack() int              { return 0 }
func (CrackedBaubleYellow) Defense() int             { return 0 }
func (CrackedBaubleYellow) Types() card.TypeSet      { return crackedBaubleTypes }
func (CrackedBaubleYellow) GoAgain() bool            { return false }

// not implemented: draft-format pitch resource; no other effect
func (CrackedBaubleYellow) NotImplemented()                              {}
func (CrackedBaubleYellow) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
