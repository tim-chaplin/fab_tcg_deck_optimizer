// Sky Fire Lanterns — Runeblade Action. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Reveal the top card of your deck. If it's <same color as this variant>, create a
// Runechant token."
//
// Peek ge.Deck[0] and compare its pitch to this variant's pitch (color). On match, create
// one Runechant.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// skyFireLanternsPlay creates a Runechant when the deck-top card matches this variant's
// pitch (color).
func skyFireLanternsPlay(ge card.GameEngine, l card.Logger, self *card.CardState, selfPitch int) {
	top, ok := ge.PeekDeck()
	if !ok || top.Pitch() != selfPitch {
		return
	}
	ge.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (c SkyFireLanternsRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	skyFireLanternsPlay(ge, l, self, c.Pitch())
}

func (c SkyFireLanternsYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	skyFireLanternsPlay(ge, l, self, c.Pitch())
}

func (c SkyFireLanternsBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	skyFireLanternsPlay(ge, l, self, c.Pitch())
}
