// Spellblade Strike — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Create a Runechant token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (SpellbladeStrikeRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (SpellbladeStrikeYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (SpellbladeStrikeBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}
