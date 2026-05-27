// Seek Horizon — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Seek Horizon, you may put a card from your hand on top
// of your deck. If you do, Seek Horizon gains **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func seekHorizonPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if !ge.MoveFromHandToTopOfDeck(self.Card.DisplayName()) {
		return
	}
	self.GrantedGoAgain = true
	l.AppendPostTriggerf(self.Card.DisplayName(), 0, "+go again")
}

func (SeekHorizonRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	seekHorizonPlay(ge, l, self)
}

func (SeekHorizonYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	seekHorizonPlay(ge, l, self)
}

func (SeekHorizonBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	seekHorizonPlay(ge, l, self)
}
