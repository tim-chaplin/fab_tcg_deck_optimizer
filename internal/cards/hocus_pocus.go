// Hocus Pocus — Runeblade Action - Attack. Cost 0, Defense 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "When this attacks, create a Runechant token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (HocusPocusRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (HocusPocusYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (HocusPocusBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}
