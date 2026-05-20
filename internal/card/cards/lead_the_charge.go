// Lead the Charge — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next time you play an action card with cost 0 or greater this turn, gain 1 action
// point. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// leadTheChargePlay grants the deferred action point eagerly when an action card still
// follows in the chain. Every action card has cost ≥ 0, so the printed "cost 0 or greater"
// clause matches any scheduled action card. Fizzles silently when none follows — an
// unspendable point would credit phantom tempo.
func leadTheChargePlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	for _, pc := range ge.CardsRemaining() {
		if pc.Card.Types(ge).Has(card.TypeAction) {
			ge.AddActionPoints(1)
			l.AppendPostTriggerf(self.Card.DisplayName(), 0, "Granted 1 action point for %s", pc.Card.DisplayName())
			return
		}
	}
}

func (LeadTheChargeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	leadTheChargePlay(ge, l, self)
}

func (LeadTheChargeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	leadTheChargePlay(ge, l, self)
}

func (LeadTheChargeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	leadTheChargePlay(ge, l, self)
}
