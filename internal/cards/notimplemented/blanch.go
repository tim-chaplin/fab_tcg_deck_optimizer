// Blanch — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, cards they own lose all colors until the end of their next turn."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: opponent 'lose all colors' debuff

func (BlanchRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: opponent 'lose all colors' debuff

func (BlanchYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: opponent 'lose all colors' debuff

func (BlanchBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
