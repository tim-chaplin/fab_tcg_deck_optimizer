// Fyendal's Fighting Spirit — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue
// 5. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks or defends, if you have less {h} than an opposing hero, gain 1{h}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// fyendalsFightingSpiritApplyRider emits the 1{h} gain as a sub-line under self when the
// current hero opts into LowerHealthWanter.
func fyendalsFightingSpiritApplyRider(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if !sim.HeroWantsLowerHealth() {
		return
	}
	s.AddValue(1)
	l.LogRider(self, 1, "Gained 1 health (lower health than opposing hero)")
}

func (FyendalsFightingSpiritRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	fyendalsFightingSpiritApplyRider(s, l, self)
}

func (FyendalsFightingSpiritYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	fyendalsFightingSpiritApplyRider(s, l, self)
}

func (FyendalsFightingSpiritBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	fyendalsFightingSpiritApplyRider(s, l, self)
}
