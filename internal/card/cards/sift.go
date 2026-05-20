// Sift — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.
//
// Text: "Put up to 4 cards from your hand on the bottom of your deck, then draw that many cards.
// **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// siftPlay cycles up to 4 hand cards to the deck bottom, then draws that many — the printed
// "up to 4" gate, capped at hand size.
func siftPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	n := len(ge.Hand())
	if n > 4 {
		n = 4
	}
	for i := 0; i < n; i++ {
		ge.AppendToDeck(ge.PopHandAt(0))
	}
	for i := 0; i < n; i++ {
		ge.DrawOne()
	}
	if n > 0 {
		l.AppendPostTriggerf(self.Card.DisplayName(), 0, "Cycled %d cards to deck bottom and drew", n)
	}
}

func (SiftRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	siftPlay(ge, l, self)
}

func (SiftYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	siftPlay(ge, l, self)
}

func (SiftBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	siftPlay(ge, l, self)
}
