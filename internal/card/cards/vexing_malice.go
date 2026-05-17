// Vexing Malice — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 2.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Deal 2 arcane damage to target hero."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (VexingMaliceRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.DealArcaneDamage(l, self.Card.DisplayName(), 2)
}

func (VexingMaliceYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.DealArcaneDamage(l, self.Card.DisplayName(), 2)
}

func (VexingMaliceBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.DealArcaneDamage(l, self.Card.DisplayName(), 2)
}
