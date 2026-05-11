// Overload — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Dominate** If Overload hits, it gains **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func overloadPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if sim.LikelyToHit(self) {
		self.GrantedGoAgain = true
	}
}

func (OverloadRed) Dominate() {}
func (OverloadRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	overloadPlay(s, l, self)
}

func (OverloadYellow) Dominate() {}
func (OverloadYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	overloadPlay(s, l, self)
}

func (OverloadBlue) Dominate() {}
func (OverloadBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	overloadPlay(s, l, self)
}
