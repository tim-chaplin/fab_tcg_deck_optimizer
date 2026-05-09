// Aether Slash — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 1.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Aether Slash attacks, if a 'non-attack' action card was pitched to play it, deal 1
// arcane damage to any target."
//
// Reads self.PitchedToPlay (the cards the chain runner attributed to funding THIS copy's
// cost) to gate the +1 arcane rider.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (AetherSlashRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	aetherSlashApplyRider(s, self)
}

func (AetherSlashYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	aetherSlashApplyRider(s, self)
}

func (AetherSlashBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	aetherSlashApplyRider(s, self)
}

// aetherSlashApplyRider deals 1 arcane and emits the rider sub-line when a non-attack action
// is among the pitched cards the runner attributed to paying for this Aether Slash.
func aetherSlashApplyRider(s *sim.TurnState, self *sim.CardState) {
	for _, p := range self.PitchedToPlay {
		if p.Types().IsNonAttackAction() {
			s.DealArcaneDamage(self, 1)
			return
		}
	}
}
