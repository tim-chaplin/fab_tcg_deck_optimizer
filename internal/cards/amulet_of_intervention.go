// Amulet of Intervention — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Intervention: Prevent the next 1 damage that
// would be dealt to your hero this turn. Activate this ability only while your hero is the target
// of a source that would deal damage equal to or greater than your hero's {h}."

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var amuletOfInterventionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type AmuletOfInterventionBlue struct{}

func (AmuletOfInterventionBlue) ID() card.ID              { return card.AmuletOfInterventionBlue }
func (AmuletOfInterventionBlue) Name() string             { return "Amulet of Intervention" }
func (AmuletOfInterventionBlue) Cost(*card.TurnState) int { return 0 }
func (AmuletOfInterventionBlue) Pitch() int               { return 3 }
func (AmuletOfInterventionBlue) Attack() int              { return 0 }
func (AmuletOfInterventionBlue) Defense() int             { return 0 }
func (AmuletOfInterventionBlue) Types() card.TypeSet      { return amuletOfInterventionTypes }
func (AmuletOfInterventionBlue) GoAgain() bool            { return true }

// not implemented: Instant 1 damage prevention; gated on a lethal incoming source
func (AmuletOfInterventionBlue) NotImplemented()                              {}
func (AmuletOfInterventionBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
