// Reek of Corruption — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Reek of Corruption gains 'When this
// hits a hero, they discard a card.'"
//
// "This hits" reads only this card's own damage; co-firing runechants don't satisfy it.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// reekOfCorruptionApplyRider registers the on-hit discard rider when the aura
// precondition is satisfied.
func reekOfCorruptionApplyRider(s *sim.TurnState, self *sim.CardState) {
	if !s.HasPlayedOrCreatedAura() {
		return
	}
	self.RegisterOnHit(reekOfCorruptionOnHit)
}

// reekOfCorruptionOnHit fires the conditional "When this hits a hero, they discard a card"
// rider. Top-level so registration stays alloc-free.
func reekOfCorruptionOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	s.AddValue(sim.DiscardValue)
	s.LogRider(self, sim.DiscardValue, "On-hit discarded a card")
}

func (ReekOfCorruptionRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	reekOfCorruptionApplyRider(s, self)
}

func (ReekOfCorruptionYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	reekOfCorruptionApplyRider(s, self)
}

func (ReekOfCorruptionBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	reekOfCorruptionApplyRider(s, self)
}
