// Amulet of Ignition — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Ignition: The next ability you activate this
// turn costs {r} less. Activate this ability only if you haven't played a card or activated an
// ability this turn."
//
// Marked sim.Unplayable: a 0/0 Item whose only output is a 1{r} discount on the next activated
// ability — gated on no card played and no ability activated yet, and the sim doesn't track
// activated-ability cost as a search dimension. The optimizer would never pick it, so it's
// filtered from random / mutation pools.

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
