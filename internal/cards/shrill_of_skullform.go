// Shrill of Skullform — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Shrill of Skullform gains +3{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (ShrillOfSkullformRed) Play(s *sim.TurnState, self *sim.CardState) {
	shrillPlay(s, self)
}

func (ShrillOfSkullformYellow) Play(s *sim.TurnState, self *sim.CardState) {
	shrillPlay(s, self)
}

func (ShrillOfSkullformBlue) Play(s *sim.TurnState, self *sim.CardState) {
	shrillPlay(s, self)
}

// shrillPlay routes the +3{p} aura-in-play buff through self.BonusAttack so EffectiveAttack
// and LikelyToHit see the buffed power, then emits the chain step at the buffed value. No
// rider sub-line — this is a power buff, not a separable effect.
func shrillPlay(s *sim.TurnState, self *sim.CardState) {
	if s.HasPlayedOrCreatedAura() {
		self.BonusAttack += 3
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
