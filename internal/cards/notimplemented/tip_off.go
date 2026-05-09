// Tip-Off — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Instant** - Discard this: **Mark** target opposing hero."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: instant discard-to-mark activation

func (c TipOffRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

// not implemented: instant discard-to-mark activation

func (c TipOffYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

// not implemented: instant discard-to-mark activation

func (c TipOffBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
