// Lead the Charge — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next time you play an action card with cost 0 or greater this turn, gain 1 action
// point. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// leadTheChargePlay registers a one-shot trigger that grants 1 action point when the next
// action card is played this attack turn. The printed "cost 0 or greater" clause admits every
// action card, so the filter is a plain action-card test. The point fizzles harmlessly
// when no action card follows — an unspent action point credits no value.
func leadTheChargePlay(ge card.GameEngine, self *card.CardState) {
	ge.CreateTrigger(self.Card, triggertype.CardOrAbility, func(g card.GameEngine, l card.Logger, t card.EphemeralTrigger) {
		g.AddActionPoints(1)
		l.AppendPostTrigger(t.CardName(), "Gained 1 action point", 0)
	}, card.TypeSet.IsAction)
}

func (LeadTheChargeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	leadTheChargePlay(ge, self)
}

func (LeadTheChargeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	leadTheChargePlay(ge, self)
}

func (LeadTheChargeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	leadTheChargePlay(ge, self)
}
