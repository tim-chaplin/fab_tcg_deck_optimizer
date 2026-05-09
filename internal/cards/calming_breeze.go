// Calming Breeze — Generic Instant. Cost 0. Printed pitch variants: Red 1.
//
// Text: "The next 3 times you would be dealt damage this turn, prevent 1 of that damage."
// The 3-event split is dropped — prevention is credited as a flat 3 against the next hit.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (CalmingBreezeRed) DefensiveInstant() {}
func (CalmingBreezeRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}
