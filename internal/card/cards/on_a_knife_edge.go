// On a Knife Edge — Generic Action. Cost 0, Defense 2. Printed only in Yellow (pitch 2).
// Text: "Your next sword attack this turn gains **go again**. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (OnAKnifeEdgeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	for _, pc := range ge.CardsRemaining() {
		if card.IsSwordAttack(ge, pc) {
			pc.GrantedGoAgain = true
			return
		}
	}
}
