// Smash Up — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 2. Only printed in Red.
//
// Text: "When this hits a hero, turn a card in their arsenal face-up, then banish an attack action
// card from their arsenal."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: on-hit opponent-arsenal manipulation rider

func (SmashUpRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
