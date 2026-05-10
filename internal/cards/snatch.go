// Snatch — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func snatchPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	self.RegisterOnHit(snatchOnHit)
}

// snatchOnHit fires the printed "When this hits, draw a card" rider. Top-level so
// registration stays alloc-free.
func snatchOnHit(s *sim.TurnState, l sim.Logger, _ *sim.CardState, _ *sim.OnHitHandler) {
	s.DrawOne()
}

func (SnatchRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	snatchPlay(s, l, self)
}

func (SnatchYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	snatchPlay(s, l, self)
}

func (SnatchBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	snatchPlay(s, l, self)
}
