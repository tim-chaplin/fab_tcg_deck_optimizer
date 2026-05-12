// Memorial Ground — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Put target attack action card with cost 2 or less from your graveyard on top of your
// deck."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Cost predicate reads s so variable-cost targets are gated on their current cost.
func memorialGroundPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	if _, ok := g.RecycleFromGraveyardToTop(func(c card.Card) bool {
		return c.Types(nil).IsAttackAction() && c.Cost(g) <= 2
	}); ok {
		l.AppendPostTrigger(self.Card.DisplayName(), "Recycled an attack action card to top of deck", 0)
	}
}

func (MemorialGroundRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	memorialGroundPlay(g, l, self)
}

func (MemorialGroundYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	memorialGroundPlay(g, l, self)
}

func (MemorialGroundBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	memorialGroundPlay(g, l, self)
}
