// Sapwood Elixir — Generic Action. Cost 1, Pitch 1, Defense 3. Only printed in Red.
//
// Text: "Your next attack this turn gets +3{p}. You may destroy a Frailty token you control. If you
// do, gain 1{h}. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

// not implemented: Frailty health-gain rider dropped (status tokens not tracked)

func (SapwoodElixirRed) Play(s *sim.TurnState, self *sim.CardState) {
	cards.GrantNextCardBonusAttack(s, 3, cards.IsAttack)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
