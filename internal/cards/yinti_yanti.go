// Yinti Yanti: "While Yinti Yanti is attacking and you control an aura, it has +1{p}.
// While Yinti Yanti is defending and you control an aura, it has +1{d}." Both bonuses
// gate on len(s.Auras) > 0 — any aura type qualifies.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// yintiYantiBonus returns +1 when any aura is in play, else 0.
func yintiYantiBonus(s *sim.TurnState) int {
	if len(s.Auras) > 0 {
		return 1
	}
	return 0
}

func yintiYantiPlay(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += yintiYantiBonus(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

func yintiYantiBlock(s *sim.TurnState, self *sim.CardState) {
	self.BonusDefense += yintiYantiBonus(s)
}

func (YintiYantiRed) Play(s *sim.TurnState, self *sim.CardState) {
	yintiYantiPlay(s, self)
}
func (YintiYantiRed) Block(s *sim.TurnState, self *sim.CardState) {
	yintiYantiBlock(s, self)
}

func (YintiYantiYellow) Play(s *sim.TurnState, self *sim.CardState) {
	yintiYantiPlay(s, self)
}
func (YintiYantiYellow) Block(s *sim.TurnState, self *sim.CardState) {
	yintiYantiBlock(s, self)
}

func (YintiYantiBlue) Play(s *sim.TurnState, self *sim.CardState) {
	yintiYantiPlay(s, self)
}
func (YintiYantiBlue) Block(s *sim.TurnState, self *sim.CardState) {
	yintiYantiBlock(s, self)
}
