// Talisman of Warfare — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** When a source you control deals exactly 2 damage to an opposing hero, destroy
// Talisman of Warfare and all cards in all arsenals."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var talismanOfWarfareTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TalismanOfWarfareYellow struct{}

func (TalismanOfWarfareYellow) ID() ids.CardID          { return ids.TalismanOfWarfareYellow }
func (TalismanOfWarfareYellow) Name() string            { return "Talisman of Warfare" }
func (TalismanOfWarfareYellow) Cost(*sim.TurnState) int { return 0 }
func (TalismanOfWarfareYellow) Pitch() int              { return 2 }
func (TalismanOfWarfareYellow) Attack() int             { return 0 }
func (TalismanOfWarfareYellow) Defense() int            { return 0 }
func (TalismanOfWarfareYellow) Types() card.TypeSet     { return talismanOfWarfareTypes }
func (TalismanOfWarfareYellow) GoAgain() bool           { return true }

// not implemented: self-destroys + wipes all arsenals on a 2-damage hit
func (TalismanOfWarfareYellow) NotImplemented()                            {}
func (TalismanOfWarfareYellow) Play(s *sim.TurnState, self *sim.CardState) { s.LogChain(self, 0) }
