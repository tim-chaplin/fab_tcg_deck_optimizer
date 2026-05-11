// Hit the High Notes — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've played or created an aura this turn, this gets +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (HitTheHighNotesRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

func (HitTheHighNotesYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

func (HitTheHighNotesBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
func hitTheHighNotesBonus(s *sim.TurnState) int {
	if s.HasPlayedOrCreatedAura() {
		return 2
	}
	return 0
}
