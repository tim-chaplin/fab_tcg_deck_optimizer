// Trade In — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, you may discard a card. If you do, draw a card. If this was played from
// arsenal, it gains **go again**."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: discard-to-draw rider, arsenal-conditional go again

func (c TradeInRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: discard-to-draw rider, arsenal-conditional go again

func (c TradeInYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: discard-to-draw rider, arsenal-conditional go again

func (c TradeInBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
