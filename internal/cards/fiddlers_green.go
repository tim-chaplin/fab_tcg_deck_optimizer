// Fiddler's Green — Generic Block. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 1. Printed on-graveyard gain: Red 3{h}, Yellow 2{h}, Blue 1{h}.
//
// Text: "When this is put into your graveyard from anywhere, gain N{h}." (N is the printed
// variant value above.)
//
// Modelling: using this card to defend sends it to the graveyard, so the N{h} gain fires on
// the DR Play path — credited as +N damage equivalent. Pitched copies go to the bottom of the
// deck instead, so they don't trigger the rider.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// fiddlersGreenPlay emits the chain step then writes the printed N{h} as a "Gained N
// health (graveyard trigger)" sub-line under self. Health is valued 1-to-1 with damage.
func fiddlersGreenPlay(s card.GameEngine, l card.Logger, self *card.CardState, heal int) {
	s.AddValue(heal)
	l.AppendPostTriggerf(self.Card.DisplayName(), heal, "Gained %d health (graveyard trigger)", heal)
}

func (FiddlersGreenRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	fiddlersGreenPlay(s, l, self, 3)
}

func (FiddlersGreenYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	fiddlersGreenPlay(s, l, self, 2)
}

func (FiddlersGreenBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	fiddlersGreenPlay(s, l, self, 1)
}
