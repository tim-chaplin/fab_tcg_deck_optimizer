// Snatch — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func snatchPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(snatchOnHit)
}

// snatchOnHit fires the printed "When this hits, draw a card" rider. Top-level so
// registration stays alloc-free.
func snatchOnHit(g card.GameEngine, l card.Logger, _ *card.CardState, _ *card.OnHitHandler) {
	g.DrawOne()
}

func (SnatchRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	snatchPlay(g, l, self)
}

func (SnatchYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	snatchPlay(g, l, self)
}

func (SnatchBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	snatchPlay(g, l, self)
}
