// Spellblade Strike — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Create a Runechant token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (SpellbladeStrikeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	createRunechantsAndLog(ge, l, self.Card.DisplayName(), 1)
}

func (SpellbladeStrikeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	createRunechantsAndLog(ge, l, self.Card.DisplayName(), 1)
}

func (SpellbladeStrikeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	createRunechantsAndLog(ge, l, self.Card.DisplayName(), 1)
}
