// Talisman of Featherfoot — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** When an attack you control gains exactly +1{p} from an effect during the
// reaction step, destroy Talisman of Featherfoot and the attack gains **go again**."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var talismanOfFeatherfootTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TalismanOfFeatherfootYellow struct{}

func (TalismanOfFeatherfootYellow) ID() ids.CardID          { return ids.TalismanOfFeatherfootYellow }
func (TalismanOfFeatherfootYellow) Name() string            { return "Talisman of Featherfoot" }
func (TalismanOfFeatherfootYellow) Cost(*sim.TurnState) int { return 0 }
func (TalismanOfFeatherfootYellow) Pitch() int              { return 2 }
func (TalismanOfFeatherfootYellow) Attack() int             { return 0 }
func (TalismanOfFeatherfootYellow) Defense() int            { return 0 }
func (TalismanOfFeatherfootYellow) Types() card.TypeSet     { return talismanOfFeatherfootTypes }
func (TalismanOfFeatherfootYellow) GoAgain() bool           { return true }

// not implemented: self-destroys when an attack gains exactly +1{p} in the reaction step →
// grants go again
func (TalismanOfFeatherfootYellow) NotImplemented()                            {}
func (TalismanOfFeatherfootYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
