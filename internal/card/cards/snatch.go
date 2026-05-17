// Snatch — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func snatchPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(snatchOnHit)
}

// snatchOnHit fires the printed "When this hits, draw a card" rider.
func snatchOnHit(ge card.GameEngine, l card.Logger, _ *card.CardState, _ *card.OnHitHandler) {
	ge.DrawOne()
}

func (SnatchRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	snatchPlay(ge, l, self)
}

func (SnatchYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	snatchPlay(ge, l, self)
}

func (SnatchBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	snatchPlay(ge, l, self)
}
