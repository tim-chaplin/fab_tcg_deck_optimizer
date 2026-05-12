// Ravenous Rabble — Generic Action - Attack. Cost 0. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, reveal the top card of your deck. This gets -X{p}, where X is the pitch
// value of the card revealed this way. **Go again**"
//
// Peek g.Deck[0].Pitch() and subtract from base power, floored at 0. If the deck is empty, no card
// is revealed so there's no penalty.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// ravenousRabbleApplyDebuff routes the -X{p} self-debuff (X = revealed deck-top pitch)
// through self.BonusAttack. Empty deck means no penalty.
func ravenousRabbleApplyDebuff(g card.GameEngine, l card.Logger, self *card.CardState) {
	top, ok := g.PeekDeck()
	if !ok {
		return
	}
	self.BonusAttack -= top.Pitch()
}

func (RavenousRabbleRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	ravenousRabbleApplyDebuff(g, l, self)
}

func (RavenousRabbleYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	ravenousRabbleApplyDebuff(g, l, self)
}

func (RavenousRabbleBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	ravenousRabbleApplyDebuff(g, l, self)
}
