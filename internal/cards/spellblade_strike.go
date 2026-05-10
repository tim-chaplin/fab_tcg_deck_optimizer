// Spellblade Strike — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Create a Runechant token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (SpellbladeStrikeRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	s.CreateRunechants(1)
	l.LogRider(self, 1, "Created a runechant")
}

func (SpellbladeStrikeYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	s.CreateRunechants(1)
	l.LogRider(self, 1, "Created a runechant")
}

func (SpellbladeStrikeBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	s.CreateRunechants(1)
	l.LogRider(self, 1, "Created a runechant")
}
