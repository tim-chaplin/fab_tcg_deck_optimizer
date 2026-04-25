// Amulet of Echoes — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Echoes: Target hero discards 2 cards.
// Activate this ability only if they have played 2 or more cards with the same name this turn."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var amuletOfEchoesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type AmuletOfEchoesBlue struct{}

func (AmuletOfEchoesBlue) ID() card.ID                               { return card.AmuletOfEchoesBlue }
func (AmuletOfEchoesBlue) Name() string                              { return "Amulet of Echoes (Blue)" }
func (AmuletOfEchoesBlue) Cost(*card.TurnState) int                  { return 0 }
func (AmuletOfEchoesBlue) Pitch() int                                { return 3 }
func (AmuletOfEchoesBlue) Attack() int                               { return 0 }
func (AmuletOfEchoesBlue) Defense() int                              { return 0 }
func (AmuletOfEchoesBlue) Types() card.TypeSet                       { return amuletOfEchoesTypes }
func (AmuletOfEchoesBlue) GoAgain() bool                             { return true }
// not implemented: Instant 'opposing hero discards 2'; gated on a repeat-name play this turn
func (AmuletOfEchoesBlue) NotImplemented()                           {}
func (AmuletOfEchoesBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
