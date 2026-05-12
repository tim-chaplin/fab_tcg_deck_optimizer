// Arcanic Crackle — Runeblade Action - Attack. Cost 0, Defense 3, Arcane 1.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Deal 1 arcane damage to target hero."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (ArcanicCrackleRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.DealArcaneDamage(l, self.Card.DisplayName(), 1)
}

func (ArcanicCrackleYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.DealArcaneDamage(l, self.Card.DisplayName(), 1)
}

func (ArcanicCrackleBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.DealArcaneDamage(l, self.Card.DisplayName(), 1)
}
