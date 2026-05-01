// Amulet of Assertiveness — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** **Attack Reaction** - Destroy Amulet of Assertiveness: Target attack gains
// "When this hits, banish the top card of your deck. If it's an attack action card, you may play it
// this turn." Activate this ability only if you have 4 or more cards in hand."
//
// Marked sim.Unplayable: a 0/0 Item whose only output is the destroy-self attack-reaction grant
// — gated on hand ≥ 4, and the granted on-hit clause depends on opposing blocks the sim doesn't
// model. The optimizer would never pick it even with the grant approximated, so it's filtered
// from random / mutation pools.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var amuletOfAssertivenessTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type AmuletOfAssertivenessYellow struct{}

func (AmuletOfAssertivenessYellow) ID() ids.CardID          { return ids.AmuletOfAssertivenessYellow }
func (AmuletOfAssertivenessYellow) Name() string            { return "Amulet of Assertiveness" }
func (AmuletOfAssertivenessYellow) Cost(*sim.TurnState) int { return 0 }
func (AmuletOfAssertivenessYellow) Pitch() int              { return 2 }
func (AmuletOfAssertivenessYellow) Attack() int             { return 0 }
func (AmuletOfAssertivenessYellow) Defense() int            { return 0 }
func (AmuletOfAssertivenessYellow) Types() card.TypeSet     { return amuletOfAssertivenessTypes }
func (AmuletOfAssertivenessYellow) GoAgain() bool           { return true }

func (AmuletOfAssertivenessYellow) Unplayable()                                {}
func (AmuletOfAssertivenessYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
