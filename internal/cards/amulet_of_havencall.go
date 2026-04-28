// Amulet of Havencall — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** **Defense Reaction** - Destroy Amulet of Havencall: Search your deck for a
// card named Rally the Rearguard, add it to this chain link as a defending card, then shuffle.
// Activate this ability only if you have no cards in hand."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var amuletOfHavencallTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type AmuletOfHavencallBlue struct{}

func (AmuletOfHavencallBlue) ID() ids.CardID           { return ids.AmuletOfHavencallBlue }
func (AmuletOfHavencallBlue) Name() string             { return "Amulet of Havencall" }
func (AmuletOfHavencallBlue) Cost(*card.TurnState) int { return 0 }
func (AmuletOfHavencallBlue) Pitch() int               { return 3 }
func (AmuletOfHavencallBlue) Attack() int              { return 0 }
func (AmuletOfHavencallBlue) Defense() int             { return 0 }
func (AmuletOfHavencallBlue) Types() card.TypeSet      { return amuletOfHavencallTypes }
func (AmuletOfHavencallBlue) GoAgain() bool            { return true }

// not implemented: DR tutor for Rally the Rearguard; gated on empty hand
func (AmuletOfHavencallBlue) NotImplemented()                              {}
func (AmuletOfHavencallBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
