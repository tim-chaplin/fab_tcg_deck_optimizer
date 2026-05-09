// Healing Balm — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 2. Printed health-gain: Red 3{h}, Yellow 2{h}, Blue 1{h}.
//
// Text: "Gain N{h}." (N is the printed variant value above.)
//
// Modelling: health is valued 1-to-1 with damage, so Play credits +N damage-equivalent per
// variant.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// healingBalmPlay emits the chain step then writes the printed N{h} as a "Gained N health"
// sub-line under self. Health is valued 1-to-1 with damage.
func healingBalmPlay(s *sim.TurnState, self *sim.CardState, heal int) {
	s.Log(self, 0)
	s.AddValue(heal)
	s.LogRiderf(self, heal, "Gained %d health", heal)
}

func (HealingBalmRed) Play(s *sim.TurnState, self *sim.CardState) {
	healingBalmPlay(s, self, 3)
}

func (HealingBalmYellow) Play(s *sim.TurnState, self *sim.CardState) {
	healingBalmPlay(s, self, 2)
}

func (HealingBalmBlue) Play(s *sim.TurnState, self *sim.CardState) {
	healingBalmPlay(s, self, 1)
}
