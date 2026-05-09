// Infectious Host — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, if you control a Frailty token, create a Frailty token under
// their control, then repeat for Inertia and Bloodrot Pox."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: frailty/inertia/bloodrot pox tokens

func (c InfectiousHostRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

// not implemented: frailty/inertia/bloodrot pox tokens

func (c InfectiousHostYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

// not implemented: frailty/inertia/bloodrot pox tokens

func (c InfectiousHostBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
