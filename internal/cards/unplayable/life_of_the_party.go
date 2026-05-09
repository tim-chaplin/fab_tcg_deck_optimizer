// Life of the Party — Generic Action - Attack. Cost 2. Printed power: Red 4, Yellow 3,
// Blue 2. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may discard or destroy a card you control named Crazy Brew rather than pay
// Life of the Party's {r} cost. If you do, choose all modes, otherwise choose 1 at random;
// - This gets 'When this hits, gain life 2{h}.'
// - This gets +2{p}.
// - This gets **go again**."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (LifeOfThePartyRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

func (LifeOfThePartyYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

func (LifeOfThePartyBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
