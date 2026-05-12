// Outed — Generic Action - Attack. Cost 0, Pitch 1, Power 3, Defense 0. Only printed in Red.
//
// Text: "If you are **marked**, you can't play this. If the defending hero is **marked**, this gets
// +1{p}. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (OutedRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	if g.OpponentMarked() {
		self.BonusAttack++
	}
}
