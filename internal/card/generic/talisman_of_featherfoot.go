// Talisman of Featherfoot — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** When an attack you control gains exactly +1{p} from an effect during the
// reaction step, destroy Talisman of Featherfoot and the attack gains **go again**."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var talismanOfFeatherfootTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TalismanOfFeatherfootYellow struct{}

func (TalismanOfFeatherfootYellow) ID() card.ID                               { return card.TalismanOfFeatherfootYellow }
func (TalismanOfFeatherfootYellow) Name() string                              { return "Talisman of Featherfoot (Yellow)" }
func (TalismanOfFeatherfootYellow) Cost(*card.TurnState) int                  { return 0 }
func (TalismanOfFeatherfootYellow) Pitch() int                                { return 2 }
func (TalismanOfFeatherfootYellow) Attack() int                               { return 0 }
func (TalismanOfFeatherfootYellow) Defense() int                              { return 0 }
func (TalismanOfFeatherfootYellow) Types() card.TypeSet                       { return talismanOfFeatherfootTypes }
func (TalismanOfFeatherfootYellow) GoAgain() bool                             { return true }
// not implemented: passive 'first attack each turn gains evade'
func (TalismanOfFeatherfootYellow) NotImplemented()                           {}
func (TalismanOfFeatherfootYellow) Play(*card.TurnState, *card.CardState) int { return 0 }
