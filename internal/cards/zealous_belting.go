// Zealous Belting — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While there is a card in your pitch zone with {p} greater than Zealous Belting's base {p},
// Zealous Belting has **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// zealousBeltingPlay grants go again when any pitched card this turn has base power greater
// than the card's own base power, then emits the chain step.
func zealousBeltingPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	base := self.Card.Attack()
	for _, p := range g.Pitched() {
		if p.Attack() > base {
			self.GrantedGoAgain = true
			break
		}
	}
}

func (ZealousBeltingRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	zealousBeltingPlay(g, l, self)
}

func (ZealousBeltingYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	zealousBeltingPlay(g, l, self)
}

func (ZealousBeltingBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	zealousBeltingPlay(g, l, self)
}
