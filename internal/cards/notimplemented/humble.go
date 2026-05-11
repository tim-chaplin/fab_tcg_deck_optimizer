// Humble — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, they lose all hero card abilities until the end of their next
// turn."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: hero-ability suppression rider

func (HumbleRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: hero-ability suppression rider

func (HumbleYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: hero-ability suppression rider

func (HumbleBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
