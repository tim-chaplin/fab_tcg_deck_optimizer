// Read the Runes — Runeblade Action. Cost 0, Defense 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Create N Runechant tokens." (Red N=3, Yellow N=2, Blue N=1.)
//
// Play returns N and sets AuraCreated so later cards this turn see the Runechants.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (ReadTheRunesRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.CreateRunechants(3)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created 3 runechants", 3)
}

func (ReadTheRunesYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.CreateRunechants(2)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created 2 runechants", 2)
}

func (ReadTheRunesBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}
