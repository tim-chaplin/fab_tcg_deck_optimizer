// Jack Be Nimble — Generic Action - Attack. Cost 0, Pitch 1, Power 3, Defense 3. Only printed in
// Red.
//
// Text: "When this attacks, you may banish a Nimblism from your graveyard. If you do, this gets
// +1{p} and **go again**. When this hits a hero, steal an item they control until the end of this
// action phase."
//
// The on-hit item-steal rider is not modelled: it's a sideboard-time consideration against
// item-heavy matchups, not a property of the card-vs-deck Value the simulator optimises.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func jackBeNimblePlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	if _, ok := g.BanishFromGraveyard(isNimblism); ok {
		self.BonusAttack++
		self.GrantedGoAgain = true
		l.AppendPostTrigger(self.Card.DisplayName(), "Banished a Nimblism, +1{p} and go again", 1)
	}
}

func (JackBeNimbleRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	jackBeNimblePlay(g, l, self)
}
