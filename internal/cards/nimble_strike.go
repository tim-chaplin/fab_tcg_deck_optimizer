// Nimble Strike — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Nimble Strike, you may banish a card named Nimblism from
// your graveyard. If you do, Nimble Strike gain +1{p} and **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func nimbleStrikePlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	if _, ok := g.BanishFromGraveyard(isNimblism); ok {
		self.BonusAttack++
		self.GrantedGoAgain = true
		l.AppendPostTrigger(self.Card.DisplayName(), "Banished a Nimblism, +1{p} and go again", 1)
	}
}

func isNimblism(c card.Card) bool { return c.Name() == "Nimblism" }

func (NimbleStrikeRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	nimbleStrikePlay(g, l, self)
}

func (NimbleStrikeYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	nimbleStrikePlay(g, l, self)
}

func (NimbleStrikeBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	nimbleStrikePlay(g, l, self)
}
