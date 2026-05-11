// Spellblade Assault — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When you attack with Spellblade Assault, create 2 Runechant tokens."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (SpellbladeAssaultRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	s.CreateRunechants(2)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created 2 runechants", 2)
}

func (SpellbladeAssaultYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	s.CreateRunechants(2)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created 2 runechants", 2)
}

func (SpellbladeAssaultBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	s.CreateRunechants(2)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created 2 runechants", 2)
}
