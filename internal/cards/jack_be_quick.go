// Jack Be Quick — Generic Action - Attack. Cost 0, Pitch 1, Power 3, Defense 3. Only printed in
// Red.
//
// Text: "When this attacks, you may banish a Nimblism from your graveyard. If you do, this gets
// +1{p} and **go again**. When this hits a hero, {u} an ally they control, then steal it until the
// end of this action phase."
//
// The on-hit ally-steal rider is not modelled: it's a sideboard-time consideration against
// ally-heavy matchups, not a property of the card-vs-deck Value the simulator optimises.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func jackBeQuickPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if _, ok := s.BanishFromGraveyard(isNimblism); ok {
		self.BonusAttack++
		self.GrantedGoAgain = true
		self.LogRider(l, 1, "Banished a Nimblism, +1{p} and go again")
	}
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

func (JackBeQuickRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	jackBeQuickPlay(s, l, self)
}
