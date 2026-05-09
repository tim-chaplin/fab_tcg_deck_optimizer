// Fate Foreseen — Generic Defense Reaction. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "Opt 1"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func fateForeseenPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
	s.Opt(1)
}

func (FateForeseenRed) Play(s *sim.TurnState, self *sim.CardState) {
	fateForeseenPlay(s, self)
}

func (FateForeseenYellow) Play(s *sim.TurnState, self *sim.CardState) {
	fateForeseenPlay(s, self)
}

func (FateForeseenBlue) Play(s *sim.TurnState, self *sim.CardState) {
	fateForeseenPlay(s, self)
}
