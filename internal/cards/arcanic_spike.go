// Arcanic Spike — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 5, Yellow 4, Blue 3.
// Text: "If you've dealt arcane damage this turn, this gets +2{p}."
//
// Rider reads TurnState.ArcaneDamageDealt: when set at Play time, +2{p}; otherwise printed
// attack alone.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// arcaneDamageBonus is the +2{p} gained when the "dealt arcane damage this turn" clause is live.
const arcaneDamageBonus = 2

// arcanicSpikeBonus returns the +2{p} power buff when ArcaneDamageDealt is set, else 0.
func arcanicSpikeBonus(s *sim.TurnState) int {
	if s != nil && s.ArcaneDamageDealt {
		return arcaneDamageBonus
	}
	return 0
}

func (ArcanicSpikeRed) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += arcanicSpikeBonus(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

func (ArcanicSpikeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += arcanicSpikeBonus(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

func (ArcanicSpikeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += arcanicSpikeBonus(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
