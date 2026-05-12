// Reek of Corruption — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Reek of Corruption gains 'When this
// hits a hero, they discard a card.'"
//
// "This hits" reads only this card's own damage; co-firing runechants don't satisfy it.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// reekOfCorruptionApplyRider registers the on-hit discard rider when the aura
// precondition is satisfied.
func reekOfCorruptionApplyRider(g card.GameEngine, l card.Logger, self *card.CardState) {
	if !g.AuraCreated() {
		return
	}
	self.RegisterOnHit(reekOfCorruptionOnHit)
}

// reekOfCorruptionOnHit fires the conditional "When this hits a hero, they discard a card"
// rider. Top-level so registration stays alloc-free.
func reekOfCorruptionOnHit(g card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	v := g.OpponentDiscard(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit discarded a card", v)
}

func (ReekOfCorruptionRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	reekOfCorruptionApplyRider(g, l, self)
}

func (ReekOfCorruptionYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	reekOfCorruptionApplyRider(g, l, self)
}

func (ReekOfCorruptionBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	reekOfCorruptionApplyRider(g, l, self)
}
