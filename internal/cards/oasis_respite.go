// Oasis Respite — Generic Instant. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text (Red): "Prevent the next 4 damage that would be dealt to target hero this turn by
// a source of your choice. If they have less life than each other hero, they may gain
// 1{h}." Yellow caps at 3, Blue at 2. The 1{h} life-gain rider fires for heroes opting
// into sim.LowerHealthWanter via sim.HeroWantsLowerHealth — life gain is credited to
// Value the same as damage prevention, so the rider lands on top of DealEffectiveDefense.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func oasisRespitePlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	if sim.HeroWantsLowerHealth() {
		s.AddValue(1)
		n++
	}
	self.Log(l, n)
}

func (OasisRespiteRed) DefensiveInstant() {}
func (OasisRespiteRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	oasisRespitePlay(s, l, self)
}

func (OasisRespiteYellow) DefensiveInstant() {}
func (OasisRespiteYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	oasisRespitePlay(s, l, self)
}

func (OasisRespiteBlue) DefensiveInstant() {}
func (OasisRespiteBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	oasisRespitePlay(s, l, self)
}
