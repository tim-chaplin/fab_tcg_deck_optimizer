// Looking for a Scrap — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Looking for a Scrap, you may banish a card with 1{p} from
// your graveyard. When you do, this gains +1{p} and **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func lookingForAScrapPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	if _, ok := s.BanishFromGraveyard(isOnePowerCard); ok {
		self.BonusAttack++
		self.GrantedGoAgain = true
		l.AppendPostTrigger(self.Card.DisplayName(), "Banished a 1{p} card, +1{p} and go again", 1)
	}
}

// isOnePowerCard matches the printed "card with 1{p}" target — any card whose printed
// Attack value is 1. Non-attack cards have Attack() = 0, so the type predicate would be
// redundant.
func isOnePowerCard(c sim.Card) bool { return c.Attack() == 1 }

func (LookingForAScrapRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	lookingForAScrapPlay(s, l, self)
}

func (LookingForAScrapYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	lookingForAScrapPlay(s, l, self)
}

func (LookingForAScrapBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	lookingForAScrapPlay(s, l, self)
}
