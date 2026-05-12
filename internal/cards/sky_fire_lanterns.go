// Sky Fire Lanterns — Runeblade Action. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Reveal the top card of your deck. If it's <same color as this variant>, create a
// Runechant token."
//
// Peek g.Deck[0] and compare its pitch to this variant's pitch (color). On match, create
// one Runechant.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// skyFireLanternsPlay emits the chain step then writes a runechant rider sub-line under
// self when the deck-top card matches this variant's pitch (color). Reads the deck top
// via PeekDeck so the cacheable bit flips — whether the rider fires depends on shuffle
// order.
func skyFireLanternsPlay(g card.GameEngine, l card.Logger, self *card.CardState, selfPitch int) {
	top, ok := g.PeekDeck()
	if !ok || top.Pitch() != selfPitch {
		return
	}
	g.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (c SkyFireLanternsRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	skyFireLanternsPlay(g, l, self, c.Pitch())
}

func (c SkyFireLanternsYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	skyFireLanternsPlay(g, l, self, c.Pitch())
}

func (c SkyFireLanternsBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	skyFireLanternsPlay(g, l, self, c.Pitch())
}
