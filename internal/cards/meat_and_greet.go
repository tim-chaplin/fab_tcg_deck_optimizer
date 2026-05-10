// Meat and Greet — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When this hits, create a Runechant token. If you've dealt arcane damage to an opposing
// hero this turn, this gets go again."
//
// The on-hit Runechant fires next turn, so it can't satisfy this card's own go-again clause.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func meatAndGreetPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if s.ArcaneDamageDealt {
		self.GrantedGoAgain = true
	}
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	self.RegisterOnHit(meatAndGreetOnHit)
}

// meatAndGreetOnHit fires the printed "When this hits, create a Runechant token" rider.
// Top-level so registration stays alloc-free.
func meatAndGreetOnHit(s *sim.TurnState, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreateRunechants(1)
	l.LogRider(self, 1, "On-hit created a runechant")
}

func (MeatAndGreetRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	meatAndGreetPlay(s, l, self)
}

func (MeatAndGreetYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	meatAndGreetPlay(s, l, self)
}

func (MeatAndGreetBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	meatAndGreetPlay(s, l, self)
}
