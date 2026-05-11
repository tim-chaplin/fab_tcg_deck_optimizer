// Drowning Dire — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 5, Yellow 4, Blue 3.
//
// Text: "If you have played or created an aura this turn, Drowning Dire gains **dominate**.
//
// When Drowning Dire hits, you may put a 'non-attack' action card from your graveyard on the
// bottom of your deck."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func drowningDireOnHitRecycle(s *sim.TurnState, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	if _, ok := s.RecycleFromGraveyardToBottom(isNonAttackAction); ok {
		self.LogRider(l, 0, "Recycled a non-attack action card to bottom of deck")
	}
}

func drowningDirePlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if s.HasPlayedOrCreatedAura() {
		self.GrantedDominate = true
	}
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
	self.RegisterOnHit(drowningDireOnHitRecycle)
}

func (DrowningDireRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	drowningDirePlay(s, l, self)
}

func (DrowningDireYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	drowningDirePlay(s, l, self)
}

func (DrowningDireBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	drowningDirePlay(s, l, self)
}
