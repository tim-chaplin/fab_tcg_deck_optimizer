// Snatch — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func snatchPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(snatchOnHit)
}

// snatchOnHit fires the printed "When this hits, draw a card" rider. Top-level so
// registration stays alloc-free.
func snatchOnHit(s card.GameEngine, l card.Logger, _ *card.CardState, _ *card.OnHitHandler) {
	s.DrawOne()
}

func (SnatchRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	snatchPlay(s, l, self)
}

func (SnatchYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	snatchPlay(s, l, self)
}

func (SnatchBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	snatchPlay(s, l, self)
}
