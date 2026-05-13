// Nimble Strike — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Nimble Strike, you may banish a card named Nimblism from
// your graveyard. If you do, Nimble Strike gain +1{p} and **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func nimbleStrikePlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if _, ok := ge.BanishFromGraveyard(isNimblism); ok {
		self.BonusAttack++
		self.GrantedGoAgain = true
		l.AppendPostTrigger(self.Card.DisplayName(), "Banished a Nimblism, +1{p} and go again", 1)
	}
}

func isNimblism(c card.Card) bool { return c.Name() == "Nimblism" }

func (NimbleStrikeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	nimbleStrikePlay(ge, l, self)
}

func (NimbleStrikeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	nimbleStrikePlay(ge, l, self)
}

func (NimbleStrikeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	nimbleStrikePlay(ge, l, self)
}
