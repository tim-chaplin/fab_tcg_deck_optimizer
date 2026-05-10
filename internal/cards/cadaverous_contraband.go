// Cadaverous Contraband — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Cadaverous Contraband hits, you may put a 'non-attack' action card from your graveyard
// on top of your deck."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func cadaverousContrabandOnHitRecycle(s *sim.TurnState, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	if _, ok := s.RecycleFromGraveyardToTop(isNonAttackAction); ok {
		l.LogRider(self, 0, "Recycled a non-attack action card to top of deck")
	}
}

func isNonAttackAction(c sim.Card) bool { return c.Types().IsNonAttackAction() }

func cadaverousContrabandPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	self.RegisterOnHit(cadaverousContrabandOnHitRecycle)
}

func (CadaverousContrabandRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	cadaverousContrabandPlay(s, l, self)
}

func (CadaverousContrabandYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	cadaverousContrabandPlay(s, l, self)
}

func (CadaverousContrabandBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	cadaverousContrabandPlay(s, l, self)
}
