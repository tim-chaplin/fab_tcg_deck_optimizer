// Memorial Ground — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Put target attack action card with cost 2 or less from your graveyard on top of your
// deck."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Cost predicate reads g so variable-cost targets are gated on their current cost.
func memorialGroundPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if _, ok := ge.RecycleFromGraveyardToTop(func(c card.Card) bool {
		return c.Types(nil).IsAttackAction() && c.Cost() <= 2
	}); ok {
		l.AppendPostTrigger(self.Card.DisplayName(), "Recycled an attack action card to top of deck", 0)
	}
}

func (MemorialGroundRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	memorialGroundPlay(ge, l, self)
}

func (MemorialGroundYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	memorialGroundPlay(ge, l, self)
}

func (MemorialGroundBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	memorialGroundPlay(ge, l, self)
}
