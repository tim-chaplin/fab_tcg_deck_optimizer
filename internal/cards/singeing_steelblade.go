// Singeing Steelblade — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 1.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When you attack with Singeing Steelblade, deal 1 arcane damage to target hero."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (SingeingSteelbladeRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.DealArcaneDamage(l, self.Card.DisplayName(), 1)
}

func (SingeingSteelbladeYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.DealArcaneDamage(l, self.Card.DisplayName(), 1)
}

func (SingeingSteelbladeBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.DealArcaneDamage(l, self.Card.DisplayName(), 1)
}
