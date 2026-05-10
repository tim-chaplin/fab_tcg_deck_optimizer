// Surging Militia — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Surging Militia has +1{p} for each non-equipment card defending it."
//
// Conservative model: the +N{p} per-defender bonus is dropped — assume the defender uses
// zero non-equipment cards to block.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func surgingMilitiaPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}

func (SurgingMilitiaRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	surgingMilitiaPlay(s, l, self)
}

func (SurgingMilitiaYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	surgingMilitiaPlay(s, l, self)
}

func (SurgingMilitiaBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	surgingMilitiaPlay(s, l, self)
}
