// Lead the Charge — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next time you play an action card with cost 0 or greater this turn, gain 1 action
// point. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// leadTheChargePlay grants Go again to the next action card scheduled this chain. Go again is
// one action point, gained as that card resolves, so the point lands on the play of the next
// action card rather than early — matching "the next time you play an action card ... gain 1
// action point". Every action card has cost ≥ 0, so the "cost 0 or greater" clause always
// matches. Fizzles when no action card follows.
func leadTheChargePlay(ge card.GameEngine) {
	for _, pc := range ge.CardsRemaining() {
		if pc.Card.Types(nil).Has(card.TypeAction) {
			pc.GrantedGoAgain = true
			return
		}
	}
}

func (LeadTheChargeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	leadTheChargePlay(ge)
}

func (LeadTheChargeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	leadTheChargePlay(ge)
}

func (LeadTheChargeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	leadTheChargePlay(ge)
}
