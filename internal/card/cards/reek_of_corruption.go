// Reek of Corruption — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Reek of Corruption gains 'When this
// hits a hero, they discard a card.'"
//
// "This hits" reads only this card's own damage; co-firing runechants don't satisfy it.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// reekOfCorruptionApplyRider registers the on-hit discard rider when the aura
// precondition is satisfied.
func reekOfCorruptionApplyRider(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if !ge.AuraCreated() {
		return
	}
	self.RegisterOnHit(reekOfCorruptionOnHit)
}

// reekOfCorruptionOnHit fires the conditional "When this hits a hero, they discard a card"
// rider.
func reekOfCorruptionOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	v := ge.OpponentDiscard(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit discarded a card", v)
}

func (ReekOfCorruptionRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	reekOfCorruptionApplyRider(ge, l, self)
}

func (ReekOfCorruptionYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	reekOfCorruptionApplyRider(ge, l, self)
}

func (ReekOfCorruptionBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	reekOfCorruptionApplyRider(ge, l, self)
}
