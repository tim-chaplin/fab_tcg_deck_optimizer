// Amulet of Intervention — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Intervention: Prevent the next 1 damage that
// would be dealt to your hero this turn. Activate this ability only while your hero is the target
// of a source that would deal damage equal to or greater than your hero's {h}."
//
// Marked sim.Unplayable: the card itself is too weak to want in a deck. Best-case output is
// preventing 1 damage when facing lethal — a card slot for 1 point of damage prevention,
// only when you're already losing. Even fully modelled the EV is far below the cost of the
// slot; the IncomingDamage modelling caveat is a secondary concern but not the deciding
// factor.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var amuletOfInterventionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type AmuletOfInterventionBlue struct{}

func (AmuletOfInterventionBlue) ID() ids.CardID          { return ids.AmuletOfInterventionBlue }
func (AmuletOfInterventionBlue) Name() string            { return "Amulet of Intervention" }
func (AmuletOfInterventionBlue) Cost(*sim.TurnState) int { return 0 }
func (AmuletOfInterventionBlue) Pitch() int              { return 3 }
func (AmuletOfInterventionBlue) Attack() int             { return 0 }
func (AmuletOfInterventionBlue) Defense() int            { return 0 }
func (AmuletOfInterventionBlue) Types() card.TypeSet     { return amuletOfInterventionTypes }
func (AmuletOfInterventionBlue) GoAgain() bool           { return true }

func (AmuletOfInterventionBlue) Unplayable()                                {}
func (AmuletOfInterventionBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
