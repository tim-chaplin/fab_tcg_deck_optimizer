// Rise Above — Generic Defense Reaction. Cost 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "You may put a card from your hand on top of your deck rather than pay Rise Above's {r}
// cost."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: hand-as-cost alt cost not modelled; card fails when printed cost can't be paid

func (RiseAboveRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}

// not implemented: hand-as-cost alt cost not modelled; card fails when printed cost can't be paid

func (RiseAboveYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}

// not implemented: hand-as-cost alt cost not modelled; card fails when printed cost can't be paid

func (RiseAboveBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}
