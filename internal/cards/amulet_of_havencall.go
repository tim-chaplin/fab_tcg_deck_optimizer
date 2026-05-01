// Amulet of Havencall — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** **Defense Reaction** - Destroy Amulet of Havencall: Search your deck for a
// card named Rally the Rearguard, add it to this chain link as a defending card, then shuffle.
// Activate this ability only if you have no cards in hand."
//
// Marked sim.Unplayable: the card itself is too weak to want in a deck. Best-case output is a
// tutored Rally the Rearguard added as a defender (~3 block), gated on empty hand and on the
// deck actually containing Rally — niche even in dedicated builds. Even fully modelled, the
// EV doesn't beat just running Rally directly; the chain-link defender plumbing is a
// secondary modelling cost but not the deciding factor.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var amuletOfHavencallTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type AmuletOfHavencallBlue struct{}

func (AmuletOfHavencallBlue) ID() ids.CardID          { return ids.AmuletOfHavencallBlue }
func (AmuletOfHavencallBlue) Name() string            { return "Amulet of Havencall" }
func (AmuletOfHavencallBlue) Cost(*sim.TurnState) int { return 0 }
func (AmuletOfHavencallBlue) Pitch() int              { return 3 }
func (AmuletOfHavencallBlue) Attack() int             { return 0 }
func (AmuletOfHavencallBlue) Defense() int            { return 0 }
func (AmuletOfHavencallBlue) Types() card.TypeSet     { return amuletOfHavencallTypes }
func (AmuletOfHavencallBlue) GoAgain() bool           { return true }

func (AmuletOfHavencallBlue) Unplayable()                                {}
func (AmuletOfHavencallBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
