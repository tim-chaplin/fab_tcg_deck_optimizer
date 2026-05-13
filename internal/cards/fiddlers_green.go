// Fiddler's Green — Generic Block. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 1. Printed on-graveyard gain: Red 3{h}, Yellow 2{h}, Blue 1{h}.
//
// Text: "When this is put into your graveyard from anywhere, gain N{h}." (N is the printed
// variant value above.)
//
// Modelling: blocking with this card sends it to the graveyard, firing the N{h} gain
// (credited as +N damage equivalent). Pitched copies bypass the rider — they go to the
// bottom of the deck.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// fiddlersGreenPlay credits N{h} as a sub-line under self (health valued 1-to-1 with damage).
func fiddlersGreenPlay(ge card.GameEngine, l card.Logger, self *card.CardState, heal int) {
	ge.AddValue(heal)
	l.AppendPostTriggerf(self.Card.DisplayName(), heal, "Gained %d health (graveyard trigger)", heal)
}

func (FiddlersGreenRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	fiddlersGreenPlay(ge, l, self, 3)
}

func (FiddlersGreenYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	fiddlersGreenPlay(ge, l, self, 2)
}

func (FiddlersGreenBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	fiddlersGreenPlay(ge, l, self, 1)
}
