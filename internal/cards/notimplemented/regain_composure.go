// Regain Composure — Generic Action. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "Your next attack this turn gets +1{p} and "When this hits, {u} your hero." **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

// not implemented: on-hit unfreeze rider (freeze/unfreeze state not tracked)

func (RegainComposureBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	cards.GrantNextCardBonusAttack(s, 1, cards.IsAttack)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
