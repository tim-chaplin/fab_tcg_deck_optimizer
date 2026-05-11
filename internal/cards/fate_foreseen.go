// Fate Foreseen — Generic Defense Reaction. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "Opt 1"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func fateForeseenPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	self.Log(l, n)
	s.Opt(l, 1)
}

func (FateForeseenRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	fateForeseenPlay(s, l, self)
}

func (FateForeseenYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	fateForeseenPlay(s, l, self)
}

func (FateForeseenBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	fateForeseenPlay(s, l, self)
}
