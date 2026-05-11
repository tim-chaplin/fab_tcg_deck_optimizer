// Even Bigger Than That! — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue
// 3.
//
// Text: "Play Even Bigger Than That! only if you've dealt {p} this turn. **Opt 3**, then reveal the
// top card of your deck. If it has {p} greater than the amount of damage you've dealt this turn,
// create a Quicken token and draw a card."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: Opt + reveal-and-Quicken trigger; gated on damage dealt this turn

func (EvenBiggerThanThatRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.Log(l, 0)
}

// not implemented: Opt + reveal-and-Quicken trigger; gated on damage dealt this turn

func (EvenBiggerThanThatYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.Log(l, 0)
}

// not implemented: Opt + reveal-and-Quicken trigger; gated on damage dealt this turn

func (EvenBiggerThanThatBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.Log(l, 0)
}
