// Talisman of Tithes — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** If an opponent would draw 1 or more cards during your action phase, instead
// destroy Talisman of Tithes and they draw that many cards minus 1."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var talismanOfTithesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TalismanOfTithesBlue struct{}

func (TalismanOfTithesBlue) ID() card.ID                               { return card.TalismanOfTithesBlue }
func (TalismanOfTithesBlue) Name() string                              { return "Talisman of Tithes (Blue)" }
func (TalismanOfTithesBlue) Cost(*card.TurnState) int                  { return 0 }
func (TalismanOfTithesBlue) Pitch() int                                { return 3 }
func (TalismanOfTithesBlue) Attack() int                               { return 0 }
func (TalismanOfTithesBlue) Defense() int                              { return 0 }
func (TalismanOfTithesBlue) Types() card.TypeSet                       { return talismanOfTithesTypes }
func (TalismanOfTithesBlue) GoAgain() bool                             { return true }
// not implemented: self-destroys on an opposing draw during your action phase → opponent draws minus 1
func (TalismanOfTithesBlue) NotImplemented()                           {}
func (TalismanOfTithesBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
