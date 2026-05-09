// Lay Low — Generic Defense Reaction. Cost 0, Pitch 2, Defense 3. Only printed in Yellow.
// Text: "If you are marked, you can't play this. If the defending hero is marked, their next attack
// this turn gets -1{p}."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: marked-defender state not tracked; treated as always legal and the -1{p}
// attacker debuff is dropped

func (LayLowYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}
