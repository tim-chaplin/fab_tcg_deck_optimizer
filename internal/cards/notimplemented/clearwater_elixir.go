// Clearwater Elixir — Generic Action. Cost 1, Pitch 1, Defense 3. Only printed in Red.
//
// Text: "Your next attack this turn gets +3{p}. You may destroy a Bloodrot Pox token you control.
// If you do, gain 1{h}. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

// not implemented: Bloodrot Pox health-gain rider dropped (status tokens not tracked)

func (ClearwaterElixirRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	cards.GrantNextCardBonusAttack(s, 3, cards.IsAttack)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
