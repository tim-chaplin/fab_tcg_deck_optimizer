// Amulet of Ignition — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Ignition: The next ability you activate this
// turn costs {r} less. Activate this ability only if you haven't played a card or activated an
// ability this turn."
//
// Marked sim.Unplayable: the card itself is too weak to want in a deck. Best-case output is a
// 1{r} discount on the next activated ability — gated on being the first action of the turn,
// so spending one card slot to save one resource on one future ability. Even fully modelled
// the EV is below the cost of the slot it occupies; the activated-ability-cost-search
// limitation is a secondary modelling concern but not the deciding factor.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var amuletOfIgnitionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type AmuletOfIgnitionYellow struct{}

func (AmuletOfIgnitionYellow) ID() ids.CardID          { return ids.AmuletOfIgnitionYellow }
func (AmuletOfIgnitionYellow) Name() string            { return "Amulet of Ignition" }
func (AmuletOfIgnitionYellow) Cost(*sim.TurnState) int { return 0 }
func (AmuletOfIgnitionYellow) Pitch() int              { return 2 }
func (AmuletOfIgnitionYellow) Attack() int             { return 0 }
func (AmuletOfIgnitionYellow) Defense() int            { return 0 }
func (AmuletOfIgnitionYellow) Types() card.TypeSet     { return amuletOfIgnitionTypes }
func (AmuletOfIgnitionYellow) GoAgain() bool           { return true }

func (AmuletOfIgnitionYellow) Unplayable()                                {}
func (AmuletOfIgnitionYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
