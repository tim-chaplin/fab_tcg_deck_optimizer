// Rally the Coast Guard — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Once per Turn Instant** - Discard a card: This gets +3{d}. Activate this only while this
// card is defending."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: defense-time instant activated ability

func (c RallyTheCoastGuardRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

// not implemented: defense-time instant activated ability

func (c RallyTheCoastGuardYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

// not implemented: defense-time instant activated ability

func (c RallyTheCoastGuardBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
