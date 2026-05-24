// Lay Low — Generic Defense Reaction. Cost 0, Pitch 2, Defense 3. Only printed in Yellow.
//
// Text: "If you are **marked**, you can't play this. If the defending hero is **marked**, their
// next attack this turn gets -1{p}."
//
// The first clause (gate on us being marked) drops — our hero's mark state isn't tracked.
// The second clause is modelled as +1{d} (BonusDefense) when the opponent is marked: we know
// the next attack would lose 1{p} to the rider, and blocking is fungible, so -1{p} attacker
// is equivalent to +1{d} on this defender.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (LayLowYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if ge.OpponentMarked() {
		self.BonusDefense += 1
	}
}
