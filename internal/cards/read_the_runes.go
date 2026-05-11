// Read the Runes — Runeblade Action. Cost 0, Defense 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Create N Runechant tokens." (Red N=3, Yellow N=2, Blue N=1.)
//
// Play returns N and sets AuraCreated so later cards this turn see the Runechants.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (ReadTheRunesRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
	s.CreateRunechants(3)
	self.LogRider(l, 3, "Created 3 runechants")
}

func (ReadTheRunesYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
	s.CreateRunechants(2)
	self.LogRider(l, 2, "Created 2 runechants")
}

func (ReadTheRunesBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
	s.CreateRunechants(1)
	self.LogRider(l, 1, "Created a runechant")
}
