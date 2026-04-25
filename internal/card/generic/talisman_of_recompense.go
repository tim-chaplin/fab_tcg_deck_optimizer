// Talisman of Recompense — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** Whenever you pitch a card, if you would gain exactly one {r}, instead destroy
// Talisman of Recompense and gain {r}{r}{r}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var talismanOfRecompenseTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TalismanOfRecompenseYellow struct{}

func (TalismanOfRecompenseYellow) ID() card.ID                               { return card.TalismanOfRecompenseYellow }
func (TalismanOfRecompenseYellow) Name() string                              { return "Talisman of Recompense (Yellow)" }
func (TalismanOfRecompenseYellow) Cost(*card.TurnState) int                  { return 0 }
func (TalismanOfRecompenseYellow) Pitch() int                                { return 2 }
func (TalismanOfRecompenseYellow) Attack() int                               { return 0 }
func (TalismanOfRecompenseYellow) Defense() int                              { return 0 }
func (TalismanOfRecompenseYellow) Types() card.TypeSet                       { return talismanOfRecompenseTypes }
func (TalismanOfRecompenseYellow) GoAgain() bool                             { return true }
// not implemented: self-destroys on pitching a 1-resource card → gain {r}{r}{r} instead
func (TalismanOfRecompenseYellow) NotImplemented()                           {}
func (TalismanOfRecompenseYellow) Play(*card.TurnState, *card.CardState) int { return 0 }
