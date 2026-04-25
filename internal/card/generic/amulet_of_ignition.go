// Amulet of Ignition — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Ignition: The next ability you activate this
// turn costs {r} less. Activate this ability only if you haven't played a card or activated an
// ability this turn."
//
// Stub only — marked NotImplemented so the optimizer skips it. The printed effect isn't modelled;
// Play returns 0.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var amuletOfIgnitionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type AmuletOfIgnitionYellow struct{}

func (AmuletOfIgnitionYellow) ID() card.ID                               { return card.AmuletOfIgnitionYellow }
func (AmuletOfIgnitionYellow) Name() string                              { return "Amulet of Ignition (Yellow)" }
func (AmuletOfIgnitionYellow) Cost(*card.TurnState) int                  { return 0 }
func (AmuletOfIgnitionYellow) Pitch() int                                { return 2 }
func (AmuletOfIgnitionYellow) Attack() int                               { return 0 }
func (AmuletOfIgnitionYellow) Defense() int                              { return 0 }
func (AmuletOfIgnitionYellow) Types() card.TypeSet                       { return amuletOfIgnitionTypes }
func (AmuletOfIgnitionYellow) GoAgain() bool                             { return true }
func (AmuletOfIgnitionYellow) NotImplemented()                           {}
func (AmuletOfIgnitionYellow) Play(*card.TurnState, *card.CardState) int { return 0 }
