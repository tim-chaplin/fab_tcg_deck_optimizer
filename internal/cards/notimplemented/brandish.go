// Brandish — Generic Action - Attack. Cost 1. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Brandish hits, your next weapon attack this turn gains +1{p}. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: next-weapon-attack +1{p} grant (weapon chain not scanned)

func (c BrandishRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: next-weapon-attack +1{p} grant (weapon chain not scanned)

func (c BrandishYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: next-weapon-attack +1{p} grant (weapon chain not scanned)

func (c BrandishBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
