// Pick a Card, Any Card — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Look at target opponent's hand then name a card. Choose a random card from their hand and
// reveal it. If it's the named card, create a Silver token. Repeat this process thrice. **Go
// again**"
//
// Approximated as a flat one Silver per play — three guesses against a typical opponent hand
// average out to roughly one match. The opponent-hand inspection itself isn't modelled.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func pickACardAnyCardPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateSilver(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a silver token", 0)
}

func (PickACardAnyCardRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	pickACardAnyCardPlay(ge, l, self)
}

func (PickACardAnyCardYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	pickACardAnyCardPlay(ge, l, self)
}

func (PickACardAnyCardBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	pickACardAnyCardPlay(ge, l, self)
}
