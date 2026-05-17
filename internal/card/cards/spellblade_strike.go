// Spellblade Strike — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Create a Runechant token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (SpellbladeStrikeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (SpellbladeStrikeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (SpellbladeStrikeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}
