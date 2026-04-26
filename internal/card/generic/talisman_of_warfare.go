// Talisman of Warfare — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** When a source you control deals exactly 2 damage to an opposing hero, destroy
// Talisman of Warfare and all cards in all arsenals."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var talismanOfWarfareTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TalismanOfWarfareYellow struct{}

func (TalismanOfWarfareYellow) ID() card.ID                               { return card.TalismanOfWarfareYellow }
func (TalismanOfWarfareYellow) Name() string                              { return "Talisman of Warfare" }
func (TalismanOfWarfareYellow) Cost(*card.TurnState) int                  { return 0 }
func (TalismanOfWarfareYellow) Pitch() int                                { return 2 }
func (TalismanOfWarfareYellow) Attack() int                               { return 0 }
func (TalismanOfWarfareYellow) Defense() int                              { return 0 }
func (TalismanOfWarfareYellow) Types() card.TypeSet                       { return talismanOfWarfareTypes }
func (TalismanOfWarfareYellow) GoAgain() bool                             { return true }
// not implemented: self-destroys + wipes all arsenals on a 2-damage hit
func (TalismanOfWarfareYellow) NotImplemented()                           {}
func (TalismanOfWarfareYellow) Play(*card.TurnState, *card.CardState) int { return 0 }
