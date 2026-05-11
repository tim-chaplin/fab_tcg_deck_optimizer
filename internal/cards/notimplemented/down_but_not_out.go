// Down But Not Out — Generic Action - Attack. Cost 3. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.
//
// Text: "When this attacks a hero, if you have less {h} and control fewer equipment and tokens than
// them, this gets +3{p}, **overpower**, and "When this hits, create an Agility, Might, and Vigor
// token.""

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: health/equipment/token comparison, agility/might/vigor tokens, overpower

func (DownButNotOutRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: health/equipment/token comparison, agility/might/vigor tokens, overpower

func (DownButNotOutYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: health/equipment/token comparison, agility/might/vigor tokens, overpower

func (DownButNotOutBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
