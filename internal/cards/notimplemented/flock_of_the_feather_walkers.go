// Flock of the Feather Walkers — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4,
// Blue 3. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Flock of the Feather Walkers, reveal a card in your hand
// with cost 1 or less. When you attack with Flock of the Feather Walkers, create a Quicken token."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: additional reveal cost, quicken tokens

func (c FlockOfTheFeatherWalkersRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}

// not implemented: additional reveal cost, quicken tokens

func (c FlockOfTheFeatherWalkersYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}

// not implemented: additional reveal cost, quicken tokens

func (c FlockOfTheFeatherWalkersBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}
