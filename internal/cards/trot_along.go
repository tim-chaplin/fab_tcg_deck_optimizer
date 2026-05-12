// Trot Along — Generic Action. Cost 0, Pitch 3, Defense 3. Only printed in Blue.
//
// Text: "Your next attack with 3 or less base {p} this turn gets **go again**. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// trotAlongApplySideEffect grants go again to the next qualifying attack scheduled later
// this turn — attack action card OR weapon swing per the "your next attack" wording —
// gated on base power 3 or less.
func trotAlongApplySideEffect(g card.GameEngine) {
	for _, pc := range g.CardsRemaining() {
		if !pc.Card.Types(nil).IsAttack() {
			continue
		}
		if pc.Card.Attack() <= 3 {
			pc.GrantedGoAgain = true
			return
		}
	}
}

func (TrotAlongBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	trotAlongApplySideEffect(g)
}
