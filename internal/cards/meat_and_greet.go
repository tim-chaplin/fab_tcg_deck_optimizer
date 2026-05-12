// Meat and Greet — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When this hits, create a Runechant token. If you've dealt arcane damage to an opposing
// hero this turn, this gets go again."
//
// The on-hit Runechant fires next turn, so it can't satisfy this card's own go-again clause.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func meatAndGreetPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	if g.ArcaneDamageDealt() {
		self.GrantedGoAgain = true
	}
	self.RegisterOnHit(meatAndGreetOnHit)
}

// meatAndGreetOnHit fires the printed "When this hits, create a Runechant token" rider.
func meatAndGreetOnHit(g card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	g.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a runechant", 1)
}

func (MeatAndGreetRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	meatAndGreetPlay(g, l, self)
}

func (MeatAndGreetYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	meatAndGreetPlay(g, l, self)
}

func (MeatAndGreetBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	meatAndGreetPlay(g, l, self)
}
