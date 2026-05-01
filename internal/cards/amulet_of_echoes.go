// Amulet of Echoes — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Echoes: Target hero discards 2 cards.
// Activate this ability only if they have played 2 or more cards with the same name this turn."
//
// Marked sim.Unplayable: a 0/0 Item whose only output is opponent-state discard, gated on the
// opponent's hand history — neither side is modelled by the sim. The optimizer would never
// pick it, so it's filtered from random / mutation pools.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var amuletOfEchoesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type AmuletOfEchoesBlue struct{}

func (AmuletOfEchoesBlue) ID() ids.CardID          { return ids.AmuletOfEchoesBlue }
func (AmuletOfEchoesBlue) Name() string            { return "Amulet of Echoes" }
func (AmuletOfEchoesBlue) Cost(*sim.TurnState) int { return 0 }
func (AmuletOfEchoesBlue) Pitch() int              { return 3 }
func (AmuletOfEchoesBlue) Attack() int             { return 0 }
func (AmuletOfEchoesBlue) Defense() int            { return 0 }
func (AmuletOfEchoesBlue) Types() card.TypeSet     { return amuletOfEchoesTypes }
func (AmuletOfEchoesBlue) GoAgain() bool           { return true }

func (AmuletOfEchoesBlue) Unplayable()                                {}
func (AmuletOfEchoesBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
