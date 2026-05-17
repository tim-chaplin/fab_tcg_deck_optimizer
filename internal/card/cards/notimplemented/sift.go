// Sift — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.
//
// Text: "Put up to 4 cards from your hand on the bottom of your deck, then draw that many cards.
// **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// not implemented: hand cycling

func (SiftRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

// not implemented: hand cycling

func (SiftYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

// not implemented: hand cycling

func (SiftBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
