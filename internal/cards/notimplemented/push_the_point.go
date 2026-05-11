// Push the Point — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If the last attack on this combat chain hit, Push the Point gains +2{p}."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: chain-history +2{p} rider (in-chain history not readable from Play)

func (c PushThePointRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: chain-history +2{p} rider (in-chain history not readable from Play)

func (c PushThePointYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: chain-history +2{p} rider (in-chain history not readable from Play)

func (c PushThePointBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
