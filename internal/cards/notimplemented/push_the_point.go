// Push the Point — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If the last attack on this combat chain hit, Push the Point gains +2{p}."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var pushThePointTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type PushThePointRed struct{}

func (PushThePointRed) ID() ids.CardID          { return ids.PushThePointRed }
func (PushThePointRed) Name() string            { return "Push the Point" }
func (PushThePointRed) Cost(*sim.TurnState) int { return 1 }
func (PushThePointRed) Pitch() int              { return 1 }
func (PushThePointRed) Attack() int             { return 4 }
func (PushThePointRed) Defense() int            { return 2 }
func (PushThePointRed) Types() card.TypeSet     { return pushThePointTypes }
func (PushThePointRed) GoAgain() bool           { return false }

// not implemented: chain-history +2{p} rider (in-chain history not readable from Play)
func (PushThePointRed) NotImplemented() {}
func (c PushThePointRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type PushThePointYellow struct{}

func (PushThePointYellow) ID() ids.CardID          { return ids.PushThePointYellow }
func (PushThePointYellow) Name() string            { return "Push the Point" }
func (PushThePointYellow) Cost(*sim.TurnState) int { return 1 }
func (PushThePointYellow) Pitch() int              { return 2 }
func (PushThePointYellow) Attack() int             { return 3 }
func (PushThePointYellow) Defense() int            { return 2 }
func (PushThePointYellow) Types() card.TypeSet     { return pushThePointTypes }
func (PushThePointYellow) GoAgain() bool           { return false }

// not implemented: chain-history +2{p} rider (in-chain history not readable from Play)
func (PushThePointYellow) NotImplemented() {}
func (c PushThePointYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type PushThePointBlue struct{}

func (PushThePointBlue) ID() ids.CardID          { return ids.PushThePointBlue }
func (PushThePointBlue) Name() string            { return "Push the Point" }
func (PushThePointBlue) Cost(*sim.TurnState) int { return 1 }
func (PushThePointBlue) Pitch() int              { return 3 }
func (PushThePointBlue) Attack() int             { return 2 }
func (PushThePointBlue) Defense() int            { return 2 }
func (PushThePointBlue) Types() card.TypeSet     { return pushThePointTypes }
func (PushThePointBlue) GoAgain() bool           { return false }

// not implemented: chain-history +2{p} rider (in-chain history not readable from Play)
func (PushThePointBlue) NotImplemented() {}
func (c PushThePointBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
